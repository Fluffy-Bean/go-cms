function Dom() {
    let blocks = {}
    let form = []

    function toggleOptionsMenu() {
        const menu = document.getElementById('options-menu')

        if (!menu) {
            return
        }

        menu.classList.toggle('options-menu--open')
    }

    async function loadAvailableBlocks() {
        const res = await fetch('/api/v1/blocks:available', {})
        const data = await res.json()

        console.log(data)

        blocks = data
    }

    async function init() {
        await loadAvailableBlocks()

        console.log('DOM script loaded')
    }

    return {
        toggleOptionsMenu,
        init,
    }
}

const dom = Dom()

window.onload = dom.init
