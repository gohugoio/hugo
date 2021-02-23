let article = document.getElementById('prose')

if (article) {
  let codeBlocks = article.getElementsByTagName('code')
    for (let [key, codeBlock] of Object.entries(codeBlocks)){
    var widthDif = codeBlock.scrollWidth - codeBlock.clientWidth
    if (widthDif > 0)
      codeBlock.parentNode.classList.add('expand')
  }
}
