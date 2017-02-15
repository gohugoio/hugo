(function() {
  //assign all blockquote content to an html collection/variable
  var blockQuotes = document.querySelectorAll('blockquote > p');
  //create new regex test for ' - ' (with one whitespace on each side so as not to accidentally grab hyphenated words as well)
  var hyphenTest = new RegExp(/\s\-\s/);
  //iterate through all html blocks within the blockquotes html collection
  for (var i = 0; i < blockQuotes.length; i++) {
    //check for ' - ' in the blockquote's text content
    if (hyphenTest.test(blockQuotes[i].textContent)) {
      //if true, split existing inner HTML into two-part array
      //newQuoteContent === text leading up to hyphen
      var newQuoteContent = blockQuotes[i].innerHTML.split(' - ')[0];
      //newAuthorAttr === text after hyphen
      var newAuthorAttr = blockQuotes[i].innerHTML.split(' - ')[1];
      //fill blockquote paragraph with new content, but now with a <cite> wrapper around the author callout and the appropriate quotation dash.
      blockQuotes[i].innerHTML = newQuoteContent + '<cite class="blockquote-citation">&#x2015; ' + newAuthorAttr + '</cite>';
    }
  }
})();