const apiHost = "https://openchart.org" // + window.location.hostname
const baseUrl = "https://openchart.org"
const maillistArchivePage = {
    searchArchiveFromApi: async function () {
            const keywords = document.getElementById('archive-search-text').value
            axios.get(apiHost+'/api/archiveSearch/'+keywords).then(function (response) {
                console.log("result:",response.data)
                var items = response.data
                var myList = document.getElementById("archiveList");

                for (var i = 0; i < items.length; i++) {
                    var listItem = document.createElement("li");
                    const link = document.createElement("a");
                    link.href = baseUrl+'/archives/text/'+items[i];
                    link.target = "_blank";
                    link.innerText = items[i];
                    listItem.appendChild(link);
                    myList.appendChild(listItem);
                }

            })
    }
}
window.onload = function () {
    document.getElementById('archive-search-btn').addEventListener('click', maillistArchivePage.searchArchiveFromApi)
}