const page = document.querySelector(".page.upload");
const element = document.querySelector("form.upload input[type=file]");
const results = document.querySelector(".uploading-images .images");
element.addEventListener("change", () => {
    page.classList.add("submitted");
    for (let file of element.files) {
        const reader = new FileReader();
        reader.addEventListener("load", async (e) => {
            const dataUrl = e.target.result;

            const image_container = document.createElement("div");
            image_container.classList.add("detail", "uploading");

            const image_title = document.createElement("h2");
            image_title.classList.add("title", "fake-input");
            image_title.contentEditable = "true";
            image_title.setAttribute("placeholder", "Title");
            image_container.appendChild(image_title);

            const image_link = document.createElement("a");
            image_link.classList.add("image");

            const image_progress = document.createElement("div");
            image_progress.classList.add("progress");

            const image_progress_indeterminate = document.createElement("div");
            image_progress_indeterminate.classList.add("indeterminate");

            image_progress.appendChild(image_progress_indeterminate);
            image_link.appendChild(image_progress);

            const image = document.createElement("img");
            image.src = dataUrl;
            image_link.appendChild(image);

            image_container.appendChild(image_link);

            const image_description = document.createElement("p");
            image_description.classList.add("description", "fake-input");
            image_description.contentEditable = "true";
            image_description.setAttribute("placeholder", "Description");
            image_description.dataset["multiline"] = "true";
            image_container.appendChild(image_description);

            results.appendChild(image_container);
            initFakeInput();

            const data = new FormData();
            data.append("file", file, file.name);

            const response = await fetch("/upload/", {
                body: data,
                cache: 'no-cache',
                credentials: 'same-origin',
                method: 'POST',
                mode: 'cors',
                redirect: 'follow'
            });
            image_container.classList.remove("uploading");
            if (response.ok) {
                const data = await response.json();
                image_link.href = "/" + data.id + ".png";
                image.src = "/" + data.id + ".png";
            } else {
                const data = await response.json();
                const image_error = document.createElement("div");
                image_error.classList.add("alert", "error");
                image_error.innerText = JSON.stringify(data.errors);
                image_container.insertBefore(image_error, image_description);
                console.log(data);
            }
        });
        reader.readAsDataURL(file);
    }
});
