// package main

// import (
// 	"fmt"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/dialog"
// 	"fyne.io/fyne/v2/driver/desktop"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	a := app.New()
// 	w := a.NewWindow("Sélection de dossier")

// 	if desc, ok := a.(desktop.App); ok {
// 		menu := fyne.NewMenu("Fichier",
// 			fyne.NewMenuItem("Quitter", func() {
// 				a.Quit()
// 			}),
// 		)
// 		desc.SetSystemTrayMenu(menu)
// 		// a.SetIcon(fyne.NewStaticResource("icon.png", []byte{}))
// 	}

// 	label := widget.NewLabel("Aucun dossier sélectionné")

// 	openBtn := widget.NewButton("Choisir un dossier", func() {
// 		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
// 			if err != nil {
// 				label.SetText("Erreur : " + err.Error())
// 				return
// 			}
// 			if uri != nil {
// 				label.SetText("Dossier sélectionné :\n" + uri.Path())
// 				fmt.Println("Chemin :", uri.Path())
// 			} else {
// 				label.SetText("Sélection annulée")
// 			}
// 		}, w)
// 	})

// 	content := container.NewVBox(
// 		openBtn,
// 		label,
// 	)

// 	w.SetContent(content)

// 	w.Resize(fyne.NewSize(1200, 800))
// 	w.ShowAndRun()
// }

package main

import (
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	var v vue = vue{
		name: "folder selected",
		path: "/selected/folder",
		a:    a,
	}

	w := a.NewWindow("Sélection de dossier")

	// Fonction d'affichage de l'interface suivante
	showNextInterface := func(folderPath string) {
		folderName := filepath.Base(folderPath)

		title := widget.NewLabelWithStyle("Nom du dossier :", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		nameLabel := widget.NewLabel(folderName)

		// Recherche d'une image
		imagePath := findFirstImage(folderPath)
		var img fyne.CanvasObject
		if imagePath != "" {
			file, err := os.Open(imagePath)
			if err != nil {
				img = widget.NewLabel("Erreur lors de l'ouverture de l'image")
			} else {
				defer file.Close()
				decoded, _, err := image.Decode(file)
				if err != nil {
					img = widget.NewLabel("Image invalide")
				} else {
					img = canvas.NewImageFromImage(decoded)
					img.(*canvas.Image).FillMode = canvas.ImageFillContain
					img.(*canvas.Image).SetMinSize(fyne.NewSize(300, 200))
				}
			}
		} else {
			img = widget.NewLabel("Aucune image trouvée dans le dossier")
		}

		backBtn := widget.NewButton("Retour", func() {
			w.SetContent(buildInitialInterface(w, v.showNextInterface))
		})

		content := container.NewVBox(
			title,
			nameLabel,
			img,
			backBtn,
		)

		w.SetContent(container.NewCenter(content))
	}

	w.SetContent(buildInitialInterface(w, showNextInterface))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}

type vue struct {
	name string
	path string
	a    fyne.App
}

func (v *vue) showNextInterface(path string) {

	label := widget.NewLabel("Aucun dossier sélectionné")
	w := v.a.NewWindow("Sélection de dossier")
	widget.NewButton("Choisir un dossier", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				label.SetText("Erreur : " + err.Error())
				return
			}
			if uri != nil {
				label.SetText("Dossier sélectionné :\n" + uri.Path())
				// onFolderChosen(uri.Path()) // transition vers la page suivante
			} else {
				label.SetText("Sélection annulée")
			}
		}, w)
	})

}

// Interface de départ
func buildInitialInterface(w fyne.Window, onFolderChosen func(path string)) fyne.CanvasObject {
	label := widget.NewLabel("Aucun dossier sélectionné")

	openBtn := widget.NewButton("Choisir un dossier", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				label.SetText("Erreur : " + err.Error())
				return
			}
			if uri != nil {
				label.SetText("Dossier sélectionné :\n" + uri.Path())
				onFolderChosen(uri.Path()) // transition vers la page suivante
			} else {
				label.SetText("Sélection annulée")
			}
		}, w)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Bienvenue", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		openBtn,
		label,
	)
}

// Trouve la première image dans le dossier (jpg, png)
func findFirstImage(folder string) string {
	var result string
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && isImageFile(path) {
			result = path
			return fs.SkipDir // stop après la première
		}
		return nil
	})
	return result
}

// Vérifie l'extension
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
