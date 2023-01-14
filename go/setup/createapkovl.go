package setup

// func CreateAPKOVL(tw *TarWriter, settings *Settings) error {
// 	staticFiles := path.Join(settings.DirSetup, "static")

// 	staticFileInfosYml, err := ioutil.ReadFile(path.Join(staticFiles, "fileinfos.yml"))
// 	if err != nil {
// 		return err
// 	}
// 	var staticFileInfos staticFileInfos
// 	yaml.Unmarshal(staticFileInfosYml, &staticFileInfos)

// 	err = filepath.Walk(staticFiles, func(fileName string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			return nil
// 		}
// 		relativePath := path.Join(".", fileName[len(staticFiles):])
// 		if linkTarget, ok := staticFileInfos.GetLinkInfo(relativePath); ok {
// 			return tw.AddLink(relativePath, linkTarget, int64(0644))
// 		} else if fileMode, ok := staticFileInfos.GetFileMode(relativePath); ok {
// 			w, err := tw.CreateFile(relativePath, fileMode, info.Size())
// 			if err != nil {
// 				return err
// 			}
// 			f, err := os.Open(fileName)
// 			if err != nil {
// 				return err
// 			}
// 			defer f.Close()
// 			_, err = io.Copy(w, f)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	templateFiles := path.Join(settings.DirSetup, "templates")
// 	err = filepath.Walk(templateFiles, func(fileName string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			return nil
// 		}
// 		relativePath := path.Join(".", fileName[len(templateFiles):])

// 		data, err := ioutil.ReadFile(fileName)
// 		if err != nil {
// 			return err
// 		}
// 		t, err := template.New(relativePath).Parse(string(data))
// 		if err != nil {
// 			return err
// 		}

// 		var buf bytes.Buffer
// 		err = t.Execute(&buf, settings)
// 		if err != nil {
// 			return err
// 		}

// 		fmt.Println(relativePath)
// 		w, err := tw.CreateFile(relativePath, int64(info.Mode()), int64(buf.Len()))
// 		if err != nil {
// 			return err
// 		}

// 		_, err = buf.WriteTo(w)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	mypiControl := path.Join(settings.DirDist, "mypi-control", "mypi-control-linux-arm64")
// 	fmt.Println("checking for mypi-control:", mypiControl)
// 	stat, err := os.Stat(mypiControl)
// 	if err != nil {
// 		return err
// 	}
// 	w2, err := tw.CreateFile("mypi-control/bin/mypi-control", 0755, stat.Size())
// 	if err != nil {
// 		return err
// 	}
// 	f, err := os.Open(mypiControl)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	_, err = io.Copy(w2, f)

// 	return err
// }
