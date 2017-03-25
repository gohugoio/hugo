(function() {
  let allVids = document.querySelectorAll('.video-thumbnail,.icon-video-play-button.shortcode');
  if (allVids.length > 0) {
    for (var i = 0; i < allVids.length; i++) {
      allVids[i].addEventListener('click', vidSwitch, false);
    }
  }
  function vidSwitch(evt) {
    var vidItem = evt.target,
      vidParent = vidItem.parentNode,
      clickedClass = vidItem.className,
      iframe = document.createElement('iframe'),
      //assign theService to the provider added, but set to lower case to control for youtube, YouTube, etc.
      theService = vidItem.parentNode.dataset.streaming.toLowerCase(),
      theVideoId = vidItem.parentNode.dataset.videoid;
    if (theService == "youtube") {
      console.log(theVideoId);
      iframe.setAttribute('src', '//www.youtube.com/embed/' + theVideoId + '?autoplay=1&autohide=2&border=0&wmode=opaque&enablejsapi=1&controls=1&showinfo=0&rel=0&vq=hd1080');
      console.log(iframe);
    } else if (theService == "vimeo") {
      iframe.setAttribute('src', '//player.vimeo.com/video/' + theVideoId + '?autoplay=1&title=0&byline=0&portrait=0');
    } else {
      console.log("If you are getting this error in the console, it is probably a sign that the youtube or vimeo api has changed.");
    }
    //The parameters for the video embed are set to show video controls but disallow related information at the video's end.
    iframe.setAttribute('frameborder', '0');
    iframe.setAttribute('class', 'video-iframe');
    if (clickedClass === "video-thumbnail" || clickedClass === "icon-video-play-button shortcode") {
      vidParent.querySelector('.icon-video-play-button').remove();
      vidParent.querySelector('.video-thumbnail').remove();
      vidParent.appendChild(iframe);
    }
  }
})();
