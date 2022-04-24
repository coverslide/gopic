import MainApp from "./elements/MainApp.js"
import GopicMain from "./elements/GopicMain.js";
import FileEntry from "./elements/FileEntry.js";
import FileList from "./elements/FileList.js";
import FolderIcon from "./elements/FolderIcon.js";
import BreadCrumb from "./elements/BreadCrumb.js";
import BreadCrumbLink from "./elements/BreadCrumbLink.js";

function defineElement(elementClass) {
    if (!elementClass.tag) {
        throw new Error("`tag` attribute is required");
    }
    customElements.define(elementClass.tag, elementClass);
    if (elementClass.template) {
        const template = document.createElement('template');
        template.innerHTML = elementClass.template;
        elementClass.prototype.templateContent = template.content;
    }
}

defineElement(MainApp);
defineElement(BreadCrumb);
defineElement(BreadCrumbLink);
defineElement(GopicMain);
defineElement(FileList);
defineElement(FileEntry);
defineElement(FolderIcon);
