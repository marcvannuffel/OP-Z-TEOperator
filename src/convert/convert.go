package convert

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/teoperator/src/ffmpeg"
	"github.com/schollz/teoperator/src/op1"
)

func newName(fname string) (fname2 string) {
	fname2 = strings.TrimSuffix(fname, filepath.Ext(fname)) + "_patch.aif"
	if _, err := os.Stat(fname2); os.IsNotExist(err) {
		// does not exist
		return
	}
	for i := 2; i < 100; i++ {
		fname2 = strings.TrimSuffix(fname, filepath.Ext(fname)) + fmt.Sprintf("_patch%d.aif", i)
		if _, err := os.Stat(fname2); os.IsNotExist(err) {
			// does not exist
			return
		}

	}
	return
}

func ToSynth(fname string, baseFreq float64) (err error) {
	log.Debug(fname)
	finalName := newName(fname)
	synthPatch := op1.NewSynthSamplePatch(baseFreq)
	err = synthPatch.SaveSample(fname, finalName, false)
	if err == nil {
		fmt.Printf("converted %+v -> %s\n", fname, finalName)
	}
	return
}

func ToDrumSplice(fname string, slices int) (err error) {
	finalName := newName(fname)
	fname2, err := ffmpeg.ToMono(fname)
	defer os.Remove(fname2)
	if err != nil {
		return
	}
	op1data := op1.NewDrumPatch()
	if slices == 0 {
		segments, errSplit := ffmpeg.SplitOnSilence(fname2, -22, 0.2, -0.2)
		if errSplit != nil {
			err = errSplit
			return
		}
		for i, seg := range segments {
			if i < len(op1data.End)-2 {
				start := int64(math.Floor(math.Round(seg.Start*100)*441)) * op1.SAMPLECONVERSION
				end := int64(math.Floor(math.Round(seg.End*100)*441)) * op1.SAMPLECONVERSION
				if start > end {
					continue
				}
				if end > op1data.End[len(op1data.End)-1] {
					continue
				}
				op1data.Start[i] = start
				op1data.End[i] = end
			}
		}
	} else {
		var totalSamples int64
		totalSamples, _, err = ffmpeg.NumSamples(fname2)
		if err != nil {
			return
		}
		log.Debugf("found %d samples", totalSamples)
		for i := 0; i < slices; i++ {
			op1data.Start[i] = int64(i) * totalSamples / int64(slices) * op1.SAMPLECONVERSION
			op1data.End[i] = int64(i+1) * totalSamples / int64(slices) * op1.SAMPLECONVERSION
		}
	}

	err = op1data.Save(fname2, finalName)
	if err == nil {
		fmt.Printf("converted %+v -> %s\n", fname, finalName)
	}
	return
}

func ToDrum(fnames []string, slices int) (err error) {
	if len(fnames) == 0 {
		err = fmt.Errorf("no files!")
		return
	}
	if len(fnames) == 1 {
		return ToDrumSplice(fnames[0], slices)
	}
	_, finalName := filepath.Split(fnames[0])
	finalName = newName(finalName)
	log.Debugf("converting %+v", fnames)
	f, err := ioutil.TempFile(".", "concat")
	defer os.Remove(f.Name())
	sampleEnd := make([]int64, len(fnames))
	fnames2 := make([]string, len(fnames))
	for i, fname := range fnames {
		var fname2 string
		fname2, err = ffmpeg.ToMono(fname)
		defer os.Remove(fname2)
		if err != nil {
			return
		}
		_, fnames2[i] = filepath.Split(fname2)
		sampleEnd[i], _, err = ffmpeg.NumSamples(fname2)
		if err != nil {
			return
		}
		if i > 0 {
			sampleEnd[i] = sampleEnd[i] + sampleEnd[i-1]
		}
		if sampleEnd[i] > 44100*12 {
			sampleEnd[i] = 44100 * 12
		}
		sampleEnd[i] = sampleEnd[i]
		log.Debugf("%s end: %d", fname, sampleEnd[i])
	}
	f.Close()

	log.Debug(fnames)
	fname2, err := ffmpeg.Concatenate(fnames2)
	defer os.Remove(fname2)
	if err != nil {
		return
	}

	drumPatch := op1.NewDrumPatch()
	for i, _ := range drumPatch.Start {
		if i == len(sampleEnd) {
			break
		}
		if i == 0 {
			drumPatch.Start[i] = 0
		} else {
			drumPatch.Start[i] = (sampleEnd[i-1]) * op1.SAMPLECONVERSION
		}
		drumPatch.End[i] = (sampleEnd[i]) * op1.SAMPLECONVERSION
	}

	err = drumPatch.Save(fname2, finalName)
	if err == nil {
		fmt.Printf("converted %+v -> %s\n", fnames, finalName)
	}
	return
}


func ToDrum2(fnames []string, slices int) (finalName string, err error) {
	if len(fnames) == 0 {
		err = fmt.Errorf("no files!")
		return
	}
	if len(fnames) == 1 {
		err = ToDrumSplice(fnames[0], slices)
		return
	}
	_, finalName = filepath.Split(fnames[0])
	finalName = newName(finalName)
	log.Debugf("converting %+v", fnames)
	f, err := ioutil.TempFile(".", "concat")
	defer os.Remove(f.Name())
	sampleEnd := make([]int64, len(fnames))
	fnames2 := make([]string, len(fnames))
	for i, fname := range fnames {
		var fname2 string
		fname2, err = ffmpeg.ToMono(fname)
		defer os.Remove(fname2)
		if err != nil {
			return
		}
		_, fnames2[i] = filepath.Split(fname2)
		sampleEnd[i], _, err = ffmpeg.NumSamples(fname2)
		if err != nil {
			return
		}
		if i > 0 {
			sampleEnd[i] = sampleEnd[i] + sampleEnd[i-1]
		}
		if sampleEnd[i] > 44100*12 {
			sampleEnd[i] = 44100 * 12
		}
		sampleEnd[i] = sampleEnd[i]
		log.Debugf("%s end: %d", fname, sampleEnd[i])
	}
	f.Close()

	log.Debug(fnames)
	fname2, err := ffmpeg.Concatenate(fnames2)
	defer os.Remove(fname2)
	if err != nil {
		return
	}

	drumPatch := op1.NewDrumPatch()
	for i, _ := range drumPatch.Start {
		if i == len(sampleEnd) {
			break
		}
		if i == 0 {
			drumPatch.Start[i] = 0
		} else {
			drumPatch.Start[i] = (sampleEnd[i-1]) * op1.SAMPLECONVERSION
		}
		drumPatch.End[i] = (sampleEnd[i]) * op1.SAMPLECONVERSION
	}

	err = drumPatch.Save(fname2, finalName)
	if err == nil {
		fmt.Printf("converted %+v -> %s\n", fnames, finalName)
	}
	return
}
// ToDrumKits converts a list of audio files into one or more drum kit .aif files.
// Files are grouped into batches so that no single kit exceeds maxSeconds total duration.
// If compress=true, samples longer than 1s are sped up 2x (no pitch shift) before packing.
// Each kit is saved as outDir/baseName_1.aif, outDir/baseName_2.aif, etc.
// Returns the list of generated .aif file paths.
func ToDrumKits(fnames []string, maxSeconds float64, outDir string, baseName string, compress bool, forceMono bool) (outFiles []string, err error) {
	if len(fnames) == 0 {
		err = fmt.Errorf("no files provided")
		return
	}
	if maxSeconds <= 0 {
		maxSeconds = 12.0
	}
	maxSamples := int64(maxSeconds * 44100)

	// Convert all files to mono WAV first and measure durations
	type monoFile struct {
		path    string
		samples int64
		origIdx int
	}
	var monoFiles []monoFile
	for i, fname := range fnames {
		var mono string
		var errMono error
		if compress {
			mono, errMono = ffmpeg.ToWavHalfSpeed(fname, 1.0, forceMono)
		} else {
			mono, errMono = ffmpeg.ToWav(fname, forceMono)
		}
		if errMono != nil {
			log.Errorf("ToMono %s: %v", fname, errMono)
			continue
		}
		n, _, errN := ffmpeg.NumSamples(mono)
		if errN != nil {
			os.Remove(mono)
			log.Errorf("NumSamples %s: %v", fname, errN)
			continue
		}
		// cap individual sample at maxSamples
		if n > maxSamples {
			n = maxSamples
		}
		monoFiles = append(monoFiles, monoFile{path: mono, samples: n, origIdx: i})
	}
	defer func() {
		for _, mf := range monoFiles {
			os.Remove(mf.path)
		}
	}()

	if len(monoFiles) == 0 {
		err = fmt.Errorf("no files could be converted")
		return
	}

	// Group into batches where cumulative samples <= maxSamples AND files <= 24
	const maxFilesPerKit = 24
	type batch struct {
		files   []monoFile
		cumSamp []int64 // cumulative sample end per file within batch
	}
	var batches []batch
	var cur batch
	var curTotal int64

	for _, mf := range monoFiles {
		// split if either limit would be exceeded
		exceedsSec := curTotal+mf.samples > maxSamples
		exceedsCount := len(cur.files) >= maxFilesPerKit
		if (exceedsSec || exceedsCount) && len(cur.files) > 0 {
			// start new batch
			batches = append(batches, cur)
			cur = batch{}
			curTotal = 0
		}
		curTotal += mf.samples
		cur.cumSamp = append(cur.cumSamp, curTotal)
		cur.files = append(cur.files, mf)
	}
	if len(cur.files) > 0 {
		batches = append(batches, cur)
	}

	// Build one kit per batch
	for bIdx, b := range batches {
		kitNum := bIdx + 1
		kitName := fmt.Sprintf("%s_%d.aif", baseName, kitNum)
		if len(batches) == 1 {
			kitName = baseName + ".aif"
		}
		outPath := filepath.Join(outDir, kitName)

		// collect mono paths for this batch (use absolute paths)
		batchPaths := make([]string, len(b.files))
		for i, mf := range b.files {
			batchPaths[i] = mf.path
		}

		// concatenate
		concatPath, errCat := ffmpeg.Concatenate(batchPaths)
		if errCat != nil {
			log.Errorf("Concatenate batch %d: %v", kitNum, errCat)
			continue
		}
		defer os.Remove(concatPath)

		// build drum patch with correct start/end markers
		drumPatch := op1.NewDrumPatch()
		for i := range drumPatch.Start {
			if i >= len(b.cumSamp) {
				break
			}
			if i == 0 {
				drumPatch.Start[i] = 0
			} else {
				drumPatch.Start[i] = b.cumSamp[i-1] * op1.SAMPLECONVERSION
			}
			drumPatch.End[i] = b.cumSamp[i] * op1.SAMPLECONVERSION
		}

		errSave := drumPatch.Save(concatPath, outPath)
		if errSave != nil {
			log.Errorf("Save batch %d: %v", kitNum, errSave)
			continue
		}
		outFiles = append(outFiles, outPath)
		log.Infof("created kit %d: %s (%d files)", kitNum, outPath, len(b.files))
	}

	if len(outFiles) == 0 {
		err = fmt.Errorf("no kits could be created")
	}
	return
}
