package corpus

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode"

	"autoatend/internal/lexer"
)

// Validation thresholds of the pre-persist battery.
const (
	noiseTolerance = 0.015 // |realized - configured| noise rate
	overlapTopK    = 100   // most frequent tokens per category compared
	overlapMin     = 0.15  // below: categories are trivially separable
	overlapMax     = 0.80  // above: categories are nearly indistinguishable
	sampleSize     = 50    // comments reviewed per category
	sampleMinRunes = 10
	sampleMaxRunes = 800
)

// Pre-registered separability criterion (PLAN.md §3, 2026-10-12), fixed before the
// templates were changed. K1 is the lexical overlap check above, measured after
// stoplist v1. K2 bounds the share of comments whose clean base sentence, as a bag of
// stoplist-filtered tokens, also occurs in another category (controlled ambiguity).
// K3 bounds the Macro-F1 of a fixed Naive Bayes probe: distinguishable, not trivial.
const (
	collisionGlobalMin = 0.05 // K2: global share of comments in the collision subset
	collisionGlobalMax = 0.10
	collisionCatMin    = 0.03 // K2: per category
	collisionRecMax    = 0.05 // K2: Reconocimientos (minority class)
	probeAlpha         = 1.0  // K3: Laplace smoothing of the probe
	probeFolds         = 5    // K3: stratified folds, rng seeded with Config.Seed
	probeMin           = 0.80 // K3: below, categories are not statistically distinguishable
	probeMax           = 0.99 // K3: above, the corpus is trivially separable
)

// Check is one named, pass/fail validation with a human-readable detail.
type Check struct {
	Name   string
	Pass   bool
	Detail string
}

// Validation is the outcome of the battery plus the measurements the report prints.
type Validation struct {
	Checks        []Check
	Counts        [5]int
	Targets       [5]int
	RealizedNoise float64
	DuplicateRate float64
	Length        struct {
		Min, Max int // in runes
		Mean     float64
	}
	Overlap [5][5]float64 // Jaccard of the top-K token sets (after stoplist v1)
	// K2: share of comments, per category and global, whose base bag occurs in another category.
	Collision       [5]float64
	CollisionGlobal float64
	ProbeF1         float64     // K3: mean Macro-F1 of the Naive Bayes probe over the folds
	Samples         [5][]string // sampleSize comments per category, chronological
}

// OK reports whether every check passed.
func (v Validation) OK() bool {
	for _, c := range v.Checks {
		if !c.Pass {
			return false
		}
	}
	return true
}

// Validate runs the battery of §6.1.4: per-category counts against target,
// no empty text or unfilled slot, well-formed chronological dates, calibrated
// noise rate, bounded inter-category lexical overlap and a 50-comment sample
// per category. The corpus is persisted only if every check passes.
func Validate(cs []Comment, cfg Config) Validation {
	v := Validation{Targets: targetCounts(cfg.Total)}
	for _, c := range cs {
		if c.CategoriaID >= 1 && c.CategoriaID <= 5 {
			v.Counts[c.CategoriaID-1]++
		}
	}
	add := func(name string, pass bool, detail string) {
		v.Checks = append(v.Checks, Check{name, pass, detail})
	}

	add("category counts", v.Counts == v.Targets && len(cs) == cfg.Total,
		fmt.Sprintf("got %v, target %v, rows %d/%d", v.Counts, v.Targets, len(cs), cfg.Total))

	empty := 0
	for _, c := range cs {
		if strings.TrimSpace(c.Texto) == "" {
			empty++
		}
	}
	add("no empty texts", empty == 0, fmt.Sprintf("%d empty texts", empty))

	unfilled := 0
	for _, c := range cs {
		if strings.ContainsAny(c.Texto, "{}") {
			unfilled++
		}
	}
	add("no unfilled slots", unfilled == 0, fmt.Sprintf("%d texts with unfilled slots", unfilled))

	ok, detail := checkDates(cs, cfg)
	add("dates", ok, detail)

	noisy := 0
	for _, c := range cs {
		if c.Noisy {
			noisy++
		}
	}
	if len(cs) > 0 {
		v.RealizedNoise = float64(noisy) / float64(len(cs))
	}
	add("noise rate", math.Abs(v.RealizedNoise-cfg.Noise.Rate) <= noiseTolerance,
		fmt.Sprintf("realized %.4f, configured %.2f, tolerance ±%.3f", v.RealizedNoise, cfg.Noise.Rate, noiseTolerance))

	v.Overlap = overlapMatrix(cs)
	lo, hi := 1.0, 0.0
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 5; j++ {
			lo = math.Min(lo, v.Overlap[i][j])
			hi = math.Max(hi, v.Overlap[i][j])
		}
	}
	add("lexical overlap", lo >= overlapMin && hi <= overlapMax,
		fmt.Sprintf("K1: pairwise Jaccard of top-%d tokens after stoplist v1 in [%.3f, %.3f], required within [%.2f, %.2f]", overlapTopK, lo, hi, overlapMin, overlapMax))

	v.Collision, v.CollisionGlobal = collisionRates(cs)
	ok = v.CollisionGlobal >= collisionGlobalMin && v.CollisionGlobal <= collisionGlobalMax && v.Collision[4] <= collisionRecMax
	for _, r := range v.Collision {
		ok = ok && r >= collisionCatMin
	}
	add("base collision", ok, fmt.Sprintf("K2: global %.4f (required [%.2f, %.2f]), per category %.4f / %.4f / %.4f / %.4f / %.4f (each >= %.2f, Reconocimientos <= %.2f)",
		v.CollisionGlobal, collisionGlobalMin, collisionGlobalMax, v.Collision[0], v.Collision[1], v.Collision[2], v.Collision[3], v.Collision[4], collisionCatMin, collisionRecMax))

	v.ProbeF1 = probeMacroF1(cs, cfg.Seed)
	add("probe classifier", v.ProbeF1 >= probeMin && v.ProbeF1 <= probeMax,
		fmt.Sprintf("K3: Naive Bayes probe (alpha %.0f, %d stratified folds) Macro-F1 %.4f, required within [%.2f, %.2f]", probeAlpha, probeFolds, v.ProbeF1, probeMin, probeMax))

	v.Samples = drawSamples(cs, cfg.Seed)
	bad := 0
	for _, s := range v.Samples {
		for _, t := range s {
			if !sampleLooksSane(t) {
				bad++
			}
		}
	}
	add("sample review", bad == 0, fmt.Sprintf("%d of %d sampled comments rejected (empty, unfilled slot, <%d or >%d runes, <2 words)",
		bad, sampleSize*5, sampleMinRunes, sampleMaxRunes))

	distinct := map[string]struct{}{}
	totalRunes := 0
	for i, c := range cs {
		distinct[c.Texto] = struct{}{}
		n := len([]rune(c.Texto))
		totalRunes += n
		if i == 0 || n < v.Length.Min {
			v.Length.Min = n
		}
		if n > v.Length.Max {
			v.Length.Max = n
		}
	}
	if len(cs) > 0 {
		v.DuplicateRate = 1 - float64(len(distinct))/float64(len(cs))
		v.Length.Mean = float64(totalRunes) / float64(len(cs))
	}
	return v
}

func checkDates(cs []Comment, cfg Config) (bool, string) {
	perDay := map[string]int{}
	prev := ""
	for i, c := range cs {
		if _, err := time.Parse(dateLayout, c.Fecha); err != nil {
			return false, fmt.Sprintf("row %d: malformed date %q", i, c.Fecha)
		}
		if c.Fecha < prev {
			return false, fmt.Sprintf("row %d: %s precedes %s", i, c.Fecha, prev)
		}
		prev = c.Fecha
		perDay[c.Fecha]++
	}
	days := make([]string, 0, len(perDay))
	for d := range perDay {
		days = append(days, d)
	}
	sort.Strings(days)
	for i, d := range days {
		last := i == len(days)-1
		if perDay[d] > cfg.PerDay || (!last && perDay[d] != cfg.PerDay) {
			return false, fmt.Sprintf("day %s has %d comments, want %d", d, perDay[d], cfg.PerDay)
		}
	}
	if len(days) == 0 {
		return false, "no dates"
	}
	return true, fmt.Sprintf("%d days, %s to %s, %d per day", len(days), days[0], days[len(days)-1], cfg.PerDay)
}

// tokens lowercases and splits on anything that is not a letter; digits are dropped.
func tokens(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) })
}

// contentTokens are tokens without the stoplist v1 entries: the tokens K1, K2 and K3 measure.
func contentTokens(s string) []string {
	var out []string
	for _, t := range tokens(s) {
		if _, stop := lexer.Lookup(t); !stop {
			out = append(out, t)
		}
	}
	return out
}

// bagKey is the order-free identity of a bag of content tokens: its sorted distinct tokens.
func bagKey(s string) string {
	ts := contentTokens(s)
	sort.Strings(ts)
	out := ts[:0]
	for i, t := range ts {
		if i == 0 || t != ts[i-1] {
			out = append(out, t)
		}
	}
	return strings.Join(out, " ")
}

// collisionRates measures K2: the share of comments whose clean base sentence, as a bag
// of content tokens, also occurs as the base of a comment of another category.
func collisionRates(cs []Comment) (perCat [5]float64, global float64) {
	cats := map[string]uint8{} // bag -> bitmask of categories
	keys := make([]string, len(cs))
	for i, c := range cs {
		if c.CategoriaID < 1 || c.CategoriaID > 5 {
			continue
		}
		keys[i] = bagKey(c.Base)
		cats[keys[i]] |= 1 << (c.CategoriaID - 1)
	}
	var hit, n [5]int
	total := 0
	for i, c := range cs {
		if c.CategoriaID < 1 || c.CategoriaID > 5 {
			continue
		}
		k := c.CategoriaID - 1
		n[k]++
		if m := cats[keys[i]]; m&(m-1) != 0 {
			hit[k]++
			total++
		}
	}
	for k := range perCat {
		if n[k] > 0 {
			perCat[k] = float64(hit[k]) / float64(n[k])
		}
	}
	if len(cs) > 0 {
		global = float64(total) / float64(len(cs))
	}
	return perCat, global
}

// probeMacroF1 measures K3: a multinomial Naive Bayes over content-token counts
// (Laplace alpha, vocabulary of the whole corpus), trained on four folds and scored on
// the fifth, with folds stratified by category from a dedicated rng stream (seed+2).
// It returns the mean Macro-F1 of the folds.
func probeMacroF1(cs []Comment, seed int64) float64 {
	vocab := map[string]int{}
	docs := make([][]int, len(cs))
	for i, c := range cs {
		for _, t := range contentTokens(c.Texto) {
			id, ok := vocab[t]
			if !ok {
				id = len(vocab)
				vocab[t] = id
			}
			docs[i] = append(docs[i], id)
		}
	}
	fold := make([]int, len(cs))
	rng := rand.New(rand.NewSource(seed + 2))
	var byCat [5][]int
	for i, c := range cs {
		if c.CategoriaID >= 1 && c.CategoriaID <= 5 {
			byCat[c.CategoriaID-1] = append(byCat[c.CategoriaID-1], i)
		}
	}
	for _, idx := range byCat {
		for j, p := range rng.Perm(len(idx)) {
			fold[idx[p]] = j % probeFolds
		}
	}
	sum := 0.0
	for f := 0; f < probeFolds; f++ {
		var classDocs [5]int
		var tokTotal [5]float64
		counts := make([][5]float64, len(vocab))
		for i, c := range cs {
			if fold[i] == f || c.CategoriaID < 1 || c.CategoriaID > 5 {
				continue
			}
			k := c.CategoriaID - 1
			classDocs[k]++
			for _, id := range docs[i] {
				counts[id][k]++
				tokTotal[k]++
			}
		}
		nTrain := 0
		for _, d := range classDocs {
			nTrain += d
		}
		var gold, pred []int
		for i, c := range cs {
			if fold[i] != f || c.CategoriaID < 1 || c.CategoriaID > 5 {
				continue
			}
			best, bestScore := 0, math.Inf(-1)
			for k := 0; k < 5; k++ {
				if classDocs[k] == 0 {
					continue
				}
				score := math.Log(float64(classDocs[k]) / float64(nTrain))
				den := math.Log(tokTotal[k] + probeAlpha*float64(len(vocab)))
				for _, id := range docs[i] {
					score += math.Log(counts[id][k]+probeAlpha) - den
				}
				if score > bestScore {
					best, bestScore = k, score
				}
			}
			gold = append(gold, c.CategoriaID)
			pred = append(pred, best+1)
		}
		sum += macroF1(gold, pred)
	}
	return sum / probeFolds
}

func overlapMatrix(cs []Comment) [5][5]float64 {
	var freq [5]map[string]int
	for i := range freq {
		freq[i] = map[string]int{}
	}
	for _, c := range cs {
		if c.CategoriaID < 1 || c.CategoriaID > 5 {
			continue
		}
		for _, t := range contentTokens(c.Texto) {
			freq[c.CategoriaID-1][t]++
		}
	}
	var top [5]map[string]bool
	for i, f := range freq {
		type kv struct {
			tok string
			n   int
		}
		list := make([]kv, 0, len(f))
		for t, n := range f {
			list = append(list, kv{t, n})
		}
		sort.Slice(list, func(a, b int) bool {
			if list[a].n != list[b].n {
				return list[a].n > list[b].n
			}
			return list[a].tok < list[b].tok
		})
		top[i] = map[string]bool{}
		for k := 0; k < len(list) && k < overlapTopK; k++ {
			top[i][list[k].tok] = true
		}
	}
	var m [5][5]float64
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			inter := 0
			for t := range top[i] {
				if top[j][t] {
					inter++
				}
			}
			union := len(top[i]) + len(top[j]) - inter
			if union > 0 {
				m[i][j] = float64(inter) / float64(union)
			}
		}
	}
	return m
}

// drawSamples picks sampleSize comments per category with its own rng stream
// (seed+1), so reviewing never perturbs generation. Output is chronological.
func drawSamples(cs []Comment, seed int64) [5][]string {
	rng := rand.New(rand.NewSource(seed + 1))
	var idx [5][]int
	for i, c := range cs {
		if c.CategoriaID >= 1 && c.CategoriaID <= 5 {
			idx[c.CategoriaID-1] = append(idx[c.CategoriaID-1], i)
		}
	}
	var out [5][]string
	for k := 0; k < 5; k++ {
		perm := rng.Perm(len(idx[k]))
		n := sampleSize
		if len(perm) < n {
			n = len(perm)
		}
		chosen := make([]int, n)
		for i := 0; i < n; i++ {
			chosen[i] = idx[k][perm[i]]
		}
		sort.Ints(chosen)
		for _, i := range chosen {
			out[k] = append(out[k], cs[i].Texto)
		}
	}
	return out
}

func sampleLooksSane(t string) bool {
	n := len([]rune(t))
	return strings.TrimSpace(t) != "" && !strings.ContainsAny(t, "{}") &&
		n >= sampleMinRunes && n <= sampleMaxRunes && len(strings.Fields(t)) >= 2
}

// macroF1 is the unweighted mean F1 over the classes 1..5 present in gold or pred.
func macroF1(gold, pred []int) float64 {
	var tp, fp, fn [6]int
	for i := range gold {
		if gold[i] == pred[i] {
			tp[gold[i]]++
		} else {
			fp[pred[i]]++
			fn[gold[i]]++
		}
	}
	sum, n := 0.0, 0
	for k := 1; k <= 5; k++ {
		if tp[k]+fp[k]+fn[k] == 0 {
			continue
		}
		n++
		sum += 2 * float64(tp[k]) / float64(2*tp[k]+fp[k]+fn[k])
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}
