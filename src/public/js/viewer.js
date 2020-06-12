document.addEventListener("DOMContentLoaded", () => {

    // Check if user has saved editor options
    let editorTheme = localStorage.getItem("editorTheme");
    if(editorTheme) {
        editor.setTheme(editorTheme);
        document.getElementById("theme").value = editorTheme;

        if(localStorage.getItem("editorNavTheme") == "BRIGHT") {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-dark");
            nav.classList.add("navbar-light");
            nav.classList.remove("bg-dark");
            nav.classList.add("bg-light");
            localStorage.setItem("editorNavTheme", "BRIGHT");
        } else {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-light");
            nav.classList.add("navbar-dark");
            nav.classList.remove("bg-light");
            nav.classList.add("bg-dark");
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
            localStorage.setItem("editorNavTheme", "BRIGHT");
        } else {
            let nav = document.getElementById("gnav");
            nav.classList.remove("navbar-light");
            nav.classList.add("navbar-dark");
            nav.classList.remove("bg-light");
            nav.classList.add("bg-dark");
            localStorage.setItem("editorNavTheme", "DARK");
        }

        editor.setTheme(e.target.value);
        localStorage.setItem("editorTheme", e.target.value);
    });

    // Download Note
    $("#downloadContent").attr('href', `data:text/plain;charset=utf-8,${encodeURI(editor.getValue())}`);

    document.getElementById("editNote").addEventListener("click", () => {
        localStorage.setItem("editorContent", editor.session.getValue());
        localStorage.setItem("editorMode", editor.session.getMode().$id);
        window.location.href = "/";
    });
     
});