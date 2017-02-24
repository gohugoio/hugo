$('.download-button').on('click', function() {
  var codeblock = $(this).siblings('.code-copy-content'),
	codeText = codeblock.text(),
	fileName = codeblock.attr('id');
  var downloadFile = new Blob([codeText], { type: 'text/plain;charset=utf-8' });
        saveAs(downloadFile, fileName);
});

