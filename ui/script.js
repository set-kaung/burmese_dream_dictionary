const contents = document.querySelector(".contents");
const home = document.querySelector(".home")

let indexCache = []
let searchCache = {}


function GetSearchData(id) {
    let list = document.createElement("ol")
    contents.textContent = "";
    list.textContent = "";
    if (id in searchCache) {
        console.log("Using search cache...}")
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

GetData()
home.addEventListener("click", GetData)