const onFormSubmit = (e) => {
    e.preventDefault()
    console.log(e)

    // const elements = e.target.querySelectorAll('.form-control')

    // const data = new FormData()

    // for (var i = 0; i < elements.lenght; i++) {
    //     data.append(elements[i].name, elements[i].value)
    // }


    const data = new URLSearchParams();
    for (const pair of new FormData(e.target)) {
        data.append(pair[0], pair[1]);
    }
    fetch('/calculate', {
        method: 'post',
        body: data,
    })
        .then((response) => {
            return response.text()
        })
        .then(result => {
            const resultText = document.querySelector('.result');
            resultText.innerHTML = result
            // console.log(result)
        })

    // stop reloading page
    return false
}

const onClearForm = () => {
    // reset input
    document.getElementById("inputForm").reset();
}

window.onload = () => {
    const form = document.getElementById("inputForm")
    form.addEventListener('submit', onFormSubmit)
    document.getElementById("clearButton").addEventListener("click", onClearForm);
}
