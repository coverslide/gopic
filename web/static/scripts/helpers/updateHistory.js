function updateHistory(url) {
    window.history.pushState('', '', url);
    window.dispatchEvent(new Event("locationchange", { url }));
}

export default updateHistory;