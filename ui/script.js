const contents = document.querySelector(".contents");
const home = document.querySelector(".home");
const search_query = document.querySelector("#search_query");
const search_query_btn = document.querySelector(".search_query_btn");

let indexCache = [];
let searchCache = {};
let currentData = [];

function AppendSearchToContents() {
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

function AppendResult(element, list) {
    let li = document.createElement("li")
    li.textContent = element.Content
    list.appendChild(li)
    contents.appendChild(list)
}

function AppendIndexData(element) {
    let div = document.createElement("div")
    div.id = `${element.BlogId}`
    let p = document.createElement("p")
    p.className = "content"
    p.textContent = element.BlogTitle
    p.addEventListener("click", () => {
        GetSingleBlogData(element.BlogId)
    })
    div.appendChild(p)
    contents.appendChild(div)
}

function AddSearchToNav(state) {
    const search_box = document.querySelector(".search_box");
    if (state == "index") {
        search_box.classList.remove("hidden")
    } else {
        search_box.classList.add("hidden")
    }
}

function GetSingleBlogData(id) {
    let list = document.createElement("ol")
    contents.textContent = "";
    AddSearchToNav("none")
    // list.textContent = "";
    AppendSearchToContents();
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



function GetAllBlogs() {
    AddSearchToNav("index")
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
                    AppendIndexData(element)
                })
            }).catch(err => console.log(err))
    } else {
        console.log("Using cache...")
        indexCache.forEach((element) => {
            AppendIndexData(element)
        })
    }
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

function SearchContents() {
    if (search_query.value == "") {
        return
    }
    let query = search_query.value;
    const formData = new FormData();
    formData.append("query", query);
    fetch("http://localhost:6969/search/content", {
        method: "POST",
        body: formData,
    })
        .then((response) => response.json())
        .then((response_data) => {
            contents.textContent = ""
            const ol = document.createElement("ol")
            contents.appendChild(ol)

            response_data.data.forEach((content) => {
                if (content.includes(query)) {
                    let li = document.createElement("li");
                    li.textContent = content;
                    ol.appendChild(li);
                }
            })
        })
        .catch((error) => {
            console.log(error);
        });
}

function CleanSearchQuery() {
    if (search_query.value == "") {
        GetAllBlogs();
    }
}


GetAllBlogs()
home.addEventListener("click", GetAllBlogs)
search_query_btn.addEventListener("click", SearchContents)
search_query.addEventListener("input", CleanSearchQuery)