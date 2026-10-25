sign_in_btn.onclick = () => {
    alert("Doesn't work yet")
}

sign_up_btn.onclick = () => {
    alert("Doesn't work yet")
}

const url_submit_btn = document.getElementById("url-submit")
const url_input_field = document.getElementById("url-input")
const url_id = document.getElementById("url-id")
const url_id_error = document.getElementById("url-id-error")

function displayServiceError() {
    url_id_error.style.display = "block"
    url_id.style.display = "none"
}

function displayUrlId(id) {
    url_id.style.display = "block"
    url_id_error.style.display = "none"

    url_id.innerText = url_id.href = `${window.location.origin}/${id}`
}

url_submit_btn.onclick = async () => {
    if (!url_input_field.checkValidity() || url_input_field.value === "") {
        url_input_field.reportValidity()
        return
    }

    const url_to_short = url_input_field.value.trim()
    const resp = await fetch(`${window.location.origin}/shorten`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            url: url_to_short,
        })
    })
    .then(resp => resp.json())
    .catch(_ => displayServiceError())

    try {
        const id = resp.id
        displayUrlId(id)
    } catch (error) {
        console.log(error)
        displayServiceError()
    }
}
