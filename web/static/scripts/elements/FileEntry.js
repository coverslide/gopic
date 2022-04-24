import updateHistory from "../helpers/updateHistory.js";
import CustomElement from "./CustomElement.js";

const TEMPLATE = `<style>
img {
	height: 50px;
	width: 50px;
}
a, a:link {
    font-weight: bold;
    font-family: sans-serif;
    color: #2c8898;
    text-decoration: none;
}
a:hover {
    text-decoration: underline;
}
</style>
`;

class FileEntry extends CustomElement {
    static tag = "file-entry";
    static template = TEMPLATE;

    connectedCallback() {
        const path = this.getAttribute("path").replace(/^\/$/, '');
        const filename = this.getAttribute("filename");
        const directory = this.getAttribute("directory") === "true";
        const body = document.createElement('div');
        const link = document.createElement('a');
        link.setAttribute('href', `${path}/${filename}`);
        if (directory) {
            link.appendChild(document.createElement('folder-icon'));
            link.addEventListener('click', (event) => {
                event.preventDefault();
                updateHistory(event.currentTarget.getAttribute('href'));
            })
        } else {
            const img = document.createElement('img');
            img.setAttribute("src", `${path}/${filename}?thumbnail=true`);
            link.appendChild(img);
        }
        link.appendChild(document.createTextNode(filename));
        body.appendChild(link);
        this.shadowRoot.appendChild(body);
    }
}

export default FileEntry;
