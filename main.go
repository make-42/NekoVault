package main

import (
        "github.com/gotk3/gotk3/gtk"
)

func main() {
        // Initialize GTK without parsing any command line arguments.
        gtk.Init(nil)

        // Create a new toplevel window, set its title, and connect it to the
        // "destroy" signal to exit the GTK main loop when it is destroyed.
        win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
        check(err)
        win.SetTitle("Neko Vault")
        win.Connect("destroy", func() {
                gtk.MainQuit()
        })

        // Create a box containing the widgets of the window.
        box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
        check(err)

        // Create a new label widget to show in the window.
        hintlabel, err := gtk.LabelNew("Click the buttons below to start encrypting/decrypting files.")
        check(err)

        // Create a new label widget to show status of program.
        statuslabel, err := gtk.LabelNew("Status: Idle")
        check(err)

        // Create an entry widget to enter key.
        entry, err := gtk.EntryNew()
        check(err)

        // Make entry widget hide text and set placeholder text.
        entry.SetVisibility(false)
        entry.SetPlaceholderText("Enter a key...")

        // Create progress bar to show encryption status.
        progressbar, err := gtk.ProgressBarNew()
        check(err)

        // Create button to encrypt files
        btnenc, err := gtk.ButtonNewWithLabel("Encrypt")
        check(err)
        btnenc.Connect("clicked", func() {
                filename := gtk.OpenFileChooserNative("Choose file/folder to encrypt.", win)
                entrytext, err := entry.GetText()
                check(err)
                encryptfile(filename, entrytext, progressbar, statuslabel)
        })

        // Create button to decrypt files
        btndec, err := gtk.ButtonNewWithLabel("Decrypt")
        check(err)
        btndec.Connect("clicked", func() {
                filename := gtk.OpenFileChooserNative("Choose file to decrypt.", win)
                entrytext, err := entry.GetText()
                check(err)
                decryptfile(filename, entrytext, progressbar, statuslabel)
        })

        // Add the box to the window.
        win.Add(box)

        // Start adding widgets to the box while have the hint label as a parent.
        box.PackStart(hintlabel, true, true, 0)

        // Add the status label widget.
        box.Add(statuslabel)

        // Add the entry widget.
        box.Add(entry)

        // Add the encrypt button to the box.
        box.Add(btnenc)

        // Add the decrypt button to the box.
        box.Add(btndec)

        // Add the progress bar to the box.
        box.Add(progressbar)

        // Set the default window size.
        win.SetDefaultSize(800, 600)

        // Recursively show all widgets contained in this window.
        win.ShowAll()

        // Begin executing the GTK main loop.  This blocks until
        // gtk.MainQuit() is run.
        gtk.Main()
}
