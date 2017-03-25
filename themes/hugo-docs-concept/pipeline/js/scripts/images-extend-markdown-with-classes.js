///this script *can* be used to add classes to images in markdown that append class= or id= to alt text)
//not currently being used in the Hugo docs
// (function imageClasses() {
//     // var allImgs = document.getElementsByTagName('img');
//     var allImgs = document.querySelectorAll('.body-copy img');
//     if (allImgs.length < 1) {
//         return;
//     } else {
//         applyAltClassesAndIds(allImgs);
//     }
//     function applyAltClassesAndIds(images) {
//         for (var i = 0; i < images.length; i++) {
//             if (images[i].alt.indexOf('class=') > 0) {
//                 var justText = images[i].alt.split('class=')[0];
//                 var newClass = images[i].alt.split('class=')[1];
//                 images[i].setAttribute('alt', justText);
//                 images[i].classList.add(newClass);
//             } else if (images[i].alt.indexOf('id=') > 0) {
//                 var justText = images[i].alt.split('id=')[0];
//                 var newId = images[i].alt.split('id=')[1];
//                 images[i].setAttribute('alt', justText);
//                 images[i].id = newId;
//             }
//         }
//     }
// })();