document.addEventListener("DOMContentLoaded", () => {

    // Copy our editor values into our form's text area
    let formTextItem = document.getElementById("post_content");
    formTextItem.style.display = "none";
    formTextItem.value = editor.session.getValue();

    editor.session.on("change", () => {
        formTextItem.value = editor.session.getValue();
        localStorage.setItem("editorContent", editor.session.getValue());
    });

    // Check if user has saved editor options
    let editorText = localStorage.getItem("editorContent");
    let editorMode = localStorage.getItem("editorMode");
    let editorTheme = localStorage.getItem("editorTheme");

    if(editorText) {
        editor.setValue(editorText);
    }
    if(editorMode) {
        editor.session.setMode(editorMode);
        document.getElementById("mode").value = editorMode;

        if(editorMode == "ace/mode/text") {
            editor.setOption("showLineNumbers", false);
        } else {
            editor.setOption("showLineNumbers", true);
        }

    }
    if(editorTheme) {
        editor.setTheme(editorTheme);
        document.getElementById("theme").value = editorTheme;

        if(localStorage.getItem("editorNavTheme") == "BRIGHT") {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-dark");
            nav.classList.add("navbar-light");
            nav.classList.remove("bg-dark");
            nav.classList.add("bg-light");

            let saveModal = document.getElementById("saveModalContent");
            saveModal.classList.remove("modal-dark");
            localStorage.setItem("editorNavTheme", "BRIGHT");
        } else {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-light");
            nav.classList.add("navbar-dark");
            nav.classList.remove("bg-light");
            nav.classList.add("bg-dark");

            let saveModal = document.getElementById("saveModalContent");
            saveModal.classList.add("modal-dark");
            localStorage.setItem("editorNavTheme", "DARK");
        }

    }

    // Editor theme changer
    document.getElementById("theme").addEventListener("change", (e) => {
        if(event.target.options[event.target.selectedIndex].parentElement.label == "Bright") {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-dark");
            nav.classList.add("navbar-light");
            nav.classList.remove("bg-dark");
            nav.classList.add("bg-light");

            let saveModal = document.getElementById("saveModalContent");
            saveModal.classList.remove("modal-dark");
            localStorage.setItem("editorNavTheme", "BRIGHT");
        } else {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-light");
            nav.classList.add("navbar-dark");
            nav.classList.remove("bg-light");
            nav.classList.add("bg-dark");

            let saveModal = document.getElementById("saveModalContent");
            saveModal.classList.add("modal-dark");
            localStorage.setItem("editorNavTheme", "DARK");
        }

        editor.setTheme(e.target.value);
        localStorage.setItem("editorTheme", e.target.value);
    });
     
    // Editor language changer
    document.getElementById("mode").addEventListener("change", (e) => {
        editor.session.setMode(e.target.value);
        localStorage.setItem("editorMode", e.target.value);

        if(e.target.value == "ace/mode/text") {
            editor.setOption("showLineNumbers", false);
        } else {
            editor.setOption("showLineNumbers", true);
        }
    });

    // Reset Editor
    document.getElementById("clear-content").addEventListener("click", () => {
        if(confirm("Are you sure you want to clear this note?\nThis will delete all text.")) {
            editor.setValue("");
            localStorage.setItem("editorContent", "");
        }
    });

    // Save Note
    document.getElementById("saveNote").addEventListener("click", () => {
        document.getElementById("noteForm").submit();
    });
});