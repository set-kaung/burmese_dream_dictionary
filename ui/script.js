const contents = document.querySelector(".contents");
const home = document.querySelector(".home")


let indexCache = [];
let searchCache = {};
let currentData = [];

function AppendSearch() {
    // console.log(list)
    let search_container = document.createElement("div")
    search_container.className = "search_container"
    let search = document.createElement("input")
    search.className = "search"
    search.addEventListener("input", Search)
    let button = document.createElement("button")
    button.className = "searchBtn"
    button.textContent = "Search"
    button.addEventListener("click", () => { Search(currentData) })
    search_container.appendChild(search)
    search_container.appendChild(button)
    contents.appendChild(search_container)
}

function GetSearchData(id) {
    let list = document.createElement("ol")
    contents.textContent = "";
    // list.textContent = "";
    AppendSearch(list);
    if (id in searchCache) {
        console.log("Using search cache...}")
        currentData = searchCache[id]
        searchCache[id].forEach(element => {
            AppendResult(element, list)
        });
    } else {
        console.log("Fetching new data...")
        const link = `http://localhost:6969/search?id=${id}`
        console.log(`get search called with ${link} `)
        fetch(link)
            .then((response) => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then((body) => {
                currentData = body.SearchDetails
                searchCache[id] = body.SearchDetails
                body.SearchDetails.forEach(element => {
                    AppendResult(element, list)
                });
            })
            .catch(err => console.log(err))
    }
}

function AppendResult(element, list) {
    let li = document.createElement("li")
    li.textContent = element.Content
    list.appendChild(li)
    contents.appendChild(list)
}


function GetData() {
    contents.textContent = "";
    if (indexCache.length == 0) {
        console.log("Fetching new data...")
        fetch("http://localhost:6969/")
            .then((response) => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            }).then((body) => {
                indexCache = body.Data
                body.Data.forEach((element) => {
                    AppendData(element)
                })
            }).catch(err => console.log(err))
    } else {
        console.log("Using cache...")
        indexCache.forEach((element) => {
            AppendData(element)
        })
    }
}
function AppendData(element) {
    let div = document.createElement("div")
    div.id = `${element.BlogId}`
    let p = document.createElement("p")
    p.className = "content"
    p.textContent = element.BlogTitle
    p.addEventListener("click", () => {
        GetSearchData(element.BlogId)
    })
    div.appendChild(p)
    contents.appendChild(div)
}

function Search() {
    const input = document.querySelector(".search")
    let keyword = input.value;
    let found = []
    found = currentData.filter((element) => {
        return element.Content.includes(keyword);
    })
    const ol = document.querySelector("ol")

    ol.textContent = "";
    found.forEach((ele) => {
        AppendResult(ele, ol)
    })

}
GetData()
home.addEventListener("click", GetData)


