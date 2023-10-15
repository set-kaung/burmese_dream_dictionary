"use strict";
const contents = document.querySelector(".contents");
const home = document.querySelector(".home");
const search_form = document.querySelector(".search_form")

const serverURL = "http://localhost:8081"

let currentID = 0;

let indexCache = [];
let searchCache = [];
//to hold single letter related search results
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
    li.textContent = element
    li.classList.add("fade-in")
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
    const search_form = document.querySelector(".search_form")
    if (state == "index") {
        search_form.classList.remove("hidden")
    } else {
        search_form.classList.add("hidden")
    }
}

function GetSingleBlogData(id) {
    currentID = id
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
        const link = `${serverURL}/search?id=${id}`
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



function GetAllBlogs(indexCache) {
    AddSearchToNav("index")
    contents.textContent = "";
    const link = `${serverURL}/`
    if (indexCache.length == 0) {
        console.log("Fetching new data...")
        fetch(link, { mode: 'cors' })
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
        console.log(indexCache)
        indexCache.forEach((element) => {
            AppendIndexData(element)
        })
    }
}


function Search() {
    const input = document.querySelector(".search")
    let query = input.value;
    const link = `${serverURL}/search/index?id=${currentID}&query=${query}`
    fetch(link).then((response) => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    }).then((body) => {
        console.log(body.data)
        const ol = document.querySelector("ol")
        ol.textContent = "";
        body.data.forEach((ele) => {
            console.log(ele)
            AppendResult(ele, ol)
        })
    }).catch(err => console.log(err))

}

function SearchContents(e) {
    e.preventDefault();
    const link = `${serverURL}/search/content`
    const formData = new FormData(search_form);
    fetch(link, {
        method: "POST",
        mode: 'cors',
        body: formData
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json(); // Parse the response as JSON
        })
        .then((response_data) => {
            console.log(response_data.data)
            if (response_data.data.length == 0) {
                contents.textContent = "There is nothing to see here..."
                return
            }
            contents.textContent = "";
            const ol = document.createElement("ol");
            contents.appendChild(ol);

            response_data.data.forEach((content) => {
                let li = document.createElement("li");
                li.classList.add("fade-in");
                li.textContent = content;
                ol.appendChild(li);
            });
        })
        .catch((error) => {
            console.error(error.message);
        });
}


home.addEventListener("click", () => { GetAllBlogs(indexCache) })
search_form.addEventListener("submit", (e) => { SearchContents(e) })
GetAllBlogs(indexCache)