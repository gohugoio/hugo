/****************
    Blast.js
****************/

/*! Blast.js (2.0.0): julian.com/research/blast (C) 2015 Julian Shapiro. MIT @license: en.wikipedia.org/wiki/MIT_License */

;
(function($, window, document, undefined) {

  /*********************
     Helper Functions
  *********************/

  /* IE detection. Gist: https://gist.github.com/julianshapiro/9098609 */
  var IE = (function() {
    if (document.documentMode) {
      return document.documentMode;
    } else {
      for (var i = 7; i > 0; i--) {
        var div = document.createElement("div");

        div.innerHTML = "<!--[if IE " + i + "]><span></span><![endif]-->";

        if (div.getElementsByTagName("span").length) {
          div = null;

          return i;
        }

        div = null;
      }
    }

    return undefined;
  })();

  /* Shim to prevent console.log() from throwing errors on IE<=7. */
  var console = window.console || { log: function() {}, time: function() {} };

  /*****************
      Constants
  *****************/

  var NAME = "blast",
    characterRanges = {
      latinPunctuation: "–—′’'“″„\"(«.…¡¿′’'”″“\")».…!?",
      latinLetters: "\\u0041-\\u005A\\u0061-\\u007A\\u00C0-\\u017F\\u0100-\\u01FF\\u0180-\\u027F"
    },
    Reg = {
      /* If the abbreviations RegEx is missing a title abbreviation that you find yourself needing to often escape manually, tweet me: @Shapiro. */
      abbreviations: new RegExp("[^" + characterRanges.latinLetters + "](e\\.g\\.)|(i\\.e\\.)|(mr\\.)|(mrs\\.)|(ms\\.)|(dr\\.)|(prof\\.)|(esq\\.)|(sr\\.)|(jr\\.)[^" + characterRanges.latinLetters + "]", "ig"),
      innerWordPeriod: new RegExp("[" + characterRanges.latinLetters + "]\.[" + characterRanges.latinLetters + "]", "ig"),
      onlyContainsPunctuation: new RegExp("[^" + characterRanges.latinPunctuation + "]"),
      adjoinedPunctuation: new RegExp("^[" + characterRanges.latinPunctuation + "]+|[" + characterRanges.latinPunctuation + "]+$", "g"),
      skippedElements: /(script|style|select|textarea)/i,
      hasPluginClass: new RegExp("(^| )" + NAME + "( |$)", "gi")
    };

  /****************
     $.fn.blast
  ****************/

  $.fn[NAME] = function(options) {

    /*************************
       Punctuation Escaping
    *************************/

    /* Escape likely false-positives of sentence-final periods. Escaping is performed by wrapping a character's ASCII equivalent in double curly brackets,
       which is then reversed (deencodcoded) after delimiting. */
    function encodePunctuation(text) {
      return text
        /* Escape the following Latin abbreviations and English titles: e.g., i.e., Mr., Mrs., Ms., Dr., Prof., Esq., Sr., and Jr. */
        .replace(Reg.abbreviations, function(match) {
          return match.replace(/\./g, "{{46}}");
        })
        /* Escape inner-word (non-space-delimited) periods. For example, the period inside "Blast.js". */
        .replace(Reg.innerWordPeriod, function(match) {
          return match.replace(/\./g, "{{46}}");
        });
    }

    /* Used to decode both the output of encodePunctuation() and punctuation that has been manually escaped by users. */
    function decodePunctuation(text) {
      return text.replace(/{{(\d{1,3})}}/g, function(fullMatch, subMatch) {
        return String.fromCharCode(subMatch);
      });
    }

    /******************
       DOM Traversal
    ******************/

    function wrapNode(node, opts) {
      var wrapper = document.createElement(opts.tag);

      /* Assign the element a class of "blast". */
      wrapper.className = NAME;

      /* If a custom class was provided, assign that too. */
      if (opts.customClass) {
        wrapper.className += " " + opts.customClass;

        /* If an opts.customClass is provided, generate an ID consisting of customClass and a number indicating the match's iteration. */
        if (opts.generateIndexID) {
          wrapper.id = opts.customClass + "-" + Element.blastedIndex;
        }
      }

      /* For the "all" delimiter, prevent space characters from collapsing. */
      if (opts.delimiter === "all" && /\s/.test(node.data)) {
        wrapper.style.whiteSpace = "pre-line";
      }

      /* Assign the element a class equal to its escaped inner text. Only applicable to the character and word delimiters (since they do not contain spaces). */
      if (opts.generateValueClass === true && !opts.search && (opts.delimiter === "character" || opts.delimiter === "word")) {
        var valueClass,
          text = node.data;

        /* For the word delimiter, remove adjoined punctuation, which is unlikely to be desired as part of the match -- unless the text
           consists solely of punctuation (e.g. "!!!"), in which case we leave the text as-is. */
        if (opts.delimiter === "word" && Reg.onlyContainsPunctuation.test(text)) {
          /* E: Remove punctuation that's adjoined to either side of the word match. */
          text = text.replace(Reg.adjoinedPunctuation, "");
        }

        valueClass = NAME + "-" + opts.delimiter.toLowerCase() + "-" + text.toLowerCase();

        wrapper.className += " " + valueClass;
      }

      /* Hide the wrapper elements from screenreaders now that we've set the target's aria-label attribute. */
      if (opts.aria) {
        wrapper.setAttribute("aria-hidden", "true");
      }

      wrapper.appendChild(node.cloneNode(false));

      return wrapper;
    }

    function traverseDOM(node, opts) {
      var matchPosition = -1,
        skipNodeBit = 0;

      /* Only proceed if the node is a text node and isn't empty. */
      if (node.nodeType === 3 && node.data.length) {
        /* Perform punctuation encoding/decoding once per original whole text node (before it gets split up into bits). */
        if (Element.nodeBeginning) {
          /* For the sentence delimiter, we first escape likely false-positive sentence-final punctuation. For all other delimiters,
             we must decode the user's manually-escaped punctuation so that the RegEx can match correctly (without being thrown off by characters in {{ASCII}}). */
          node.data = (!opts.search && opts.delimiter === "sentence") ? encodePunctuation(node.data) : decodePunctuation(node.data);

          Element.nodeBeginning = false;
        }

        matchPosition = node.data.search(delimiterRegex);

        /* If there's a RegEx match in this text node, proceed with element wrapping. */
        if (matchPosition !== -1) {
          var match = node.data.match(delimiterRegex),
            matchText = match[0],
            subMatchText = match[1] || false;

          /* RegEx queries that can return empty strings (e.g ".*") produce an empty matchText which throws the entire traversal process into an infinite loop due to the position index not incrementing.
             Thus, we bump up the position index manually, resulting in a zero-width split at this location followed by the continuation of the traversal process. */
          if (matchText === "") {
            matchPosition++;
            /* If a RegEx submatch is produced that is not identical to the full string match, use the submatch's index position and text.
               This technique allows us to avoid writing multi-part RegEx queries for submatch finding. */
          } else if (subMatchText && subMatchText !== matchText) {
            matchPosition += matchText.indexOf(subMatchText);
            matchText = subMatchText;
          }

          /* Split this text node into two separate nodes at the position of the match, returning the node that begins after the match position. */
          var middleBit = node.splitText(matchPosition);

          /* Split the newly-produced text node at the end of the match's text so that middleBit is a text node that consists solely of the matched text. The other newly-created text node, which begins
             at the end of the match's text, is what will be traversed in the subsequent loop (in order to find additional matches in the containing text node). */
          middleBit.splitText(matchText.length);

          /* Over-increment the loop counter (see below) so that we skip the extra node (middleBit) that we've just created (and already processed). */
          skipNodeBit = 1;

          if (!opts.search && opts.delimiter === "sentence") {
            /* Now that we've forcefully escaped all likely false-positive sentence-final punctuation, we must decode the punctuation back from ASCII. */
            middleBit.data = decodePunctuation(middleBit.data);
          }

          /* Create the wrapped node. */
          var wrappedNode = wrapNode(middleBit, opts, Element.blastedIndex);
          /* Then replace the middleBit text node with its wrapped version. */
          middleBit.parentNode.replaceChild(wrappedNode, middleBit);

          /* Push the wrapper onto the Element.wrappers array (for later use with stack manipulation). */
          Element.wrappers.push(wrappedNode);

          Element.blastedIndex++;

          /* Note: We use this slow splice-then-iterate method because every match needs to be converted into an HTML element node. A text node's text cannot have HTML elements inserted into it. */
          /* TODO: To improve performance, use documentFragments to delay node manipulation so that DOM queries and updates can be batched across elements. */
        }
        /* Traverse the DOM tree until we find text nodes. Skip script and style elements. Skip select and textarea elements since they contain special text nodes that users would not want wrapped.
           Additionally, check for the existence of our plugin's class to ensure that we do not retraverse elements that have already been blasted. */
        /* Note: This basic DOM traversal technique is copyright Johann Burkard <http://johannburkard.de>. Licensed under the MIT License: http://en.wikipedia.org/wiki/MIT_License */
      } else if (node.nodeType === 1 && node.hasChildNodes() && !Reg.skippedElements.test(node.tagName) && !Reg.hasPluginClass.test(node.className)) {
        /* Note: We don't cache childNodes' length since it's a live nodeList (which changes dynamically with the use of splitText() above). */
        for (var i = 0; i < node.childNodes.length; i++) {
          Element.nodeBeginning = true;

          i += traverseDOM(node.childNodes[i], opts);
        }
      }

      return skipNodeBit;
    }

    /*******************
       Call Variables
    *******************/

    var opts = $.extend({}, $.fn[NAME].defaults, options),
      delimiterRegex,
      /* Container for variables specific to each element targeted by the Blast call. */
      Element = {};

    /***********************
       Delimiter Creation
    ***********************/

    /* Ensure that the opts.delimiter search variable is a non-empty string. */
    if (opts.search.length && (typeof opts.search === "string" || /^\d/.test(parseFloat(opts.search)))) {
      /* Since the search is performed as a Regex (see below), we escape the string's Regex meta-characters. */
      opts.delimiter = opts.search.toString().replace(/[-[\]{,}(.)*+?|^$\\\/]/g, "\\$&");

      /* Note: This matches the apostrophe+s of the phrase's possessive form: {PHRASE's}. */
      /* Note: This will not match text that is part of a compound word (two words adjoined with a dash), e.g. "front" won't match inside "front-end". */
      /* Note: Based on the way the search algorithm is implemented, it is not possible to search for a string that consists solely of punctuation characters. */
      /* Note: By creating boundaries at Latin alphabet ranges instead of merely spaces, we effectively match phrases that are inlined alongside any type of non-Latin-letter,
         e.g. word|, word!, ♥word♥ will all match. */
      delimiterRegex = new RegExp("(?:^|[^-" + characterRanges.latinLetters + "])(" + opts.delimiter + "('s)?)(?![-" + characterRanges.latinLetters + "])", "i");
    } else {
      /* Normalize the string's case for the delimiter switch check below. */
      if (typeof opts.delimiter === "string") {
        opts.delimiter = opts.delimiter.toLowerCase();
      }

      switch (opts.delimiter) {
        case "all":
          /* Matches every character then later sets spaces to "white-space: pre-line" so they don't collapse. */
          delimiterRegex = /(.)/;
          break;

        case "letter":
        case "char":
        case "character":
          /* Matches every non-space character. */
          /* Note: This is the slowest delimiter. However, its slowness is only noticeable when it's used on larger bodies of text (of over 500 characters) on <=IE8.
             (Run Blast with opts.debug=true to monitor execution times.) */
          delimiterRegex = /(\S)/;
          break;

        case "word":
          /* Matches strings in between space characters. */
          /* Note: Matches will include any punctuation that's adjoined to the word, e.g. "Hey!" will be a full match. */
          /* Note: Remember that, with Blast, every HTML element marks the start of a brand new string. Hence, "in<b>si</b>de" matches as three separate words. */
          delimiterRegex = /\s*(\S+)\s*/;
          break;

        case "sentence":
          /* Matches phrases either ending in Latin alphabet punctuation or located at the end of the text. (Linebreaks are not considered punctuation.) */
          /* Note: If you don't want punctuation to demarcate a sentence match, replace the punctuation character with {{ASCII_CODE_FOR_DESIRED_PUNCTUATION}}. ASCII codes: .={{46}}, ?={{63}}, !={{33}} */
          delimiterRegex = /(?=\S)(([.]{2,})?[^!?]+?([.…!?]+|(?=\s+$)|$)(\s*[′’'”″“")»]+)*)/;
          /* RegExp explanation (Tip: Use Regex101.com to play around with this expression and see which strings it matches):
             - Expanded view: /(?=\S) ( ([.]{2,})? [^!?]+? ([.…!?]+|(?=\s+$)|$) (\s*[′’'”″“")»]+)* )
             - (?=\S) --> Match must contain a non-space character.
             - ([.]{2,})? --> Match may begin with a group of periods.
             - [^!?]+? --> Grab everything that isn't an unequivocally-terminating punctuation character, but stop at the following condition...
             - ([.…!?]+|(?=\s+$)|$) --> Match the last occurrence of sentence-final punctuation or the end of the text (optionally with left-side trailing spaces).
             - (\s*[′’'”″“")»]+)* --> After the final punctuation, match any and all pairs of (optionally space-delimited) quotes and parentheses.
          */
          break;

        case "element":
          /* Matches text between HTML tags. */
          /* Note: Wrapping always occurs inside of elements, i.e. <b><span class="blast">Bolded text here</span></b>. */
          delimiterRegex = /(?=\S)([\S\s]*\S)/;
          break;

          /*****************
             Custom Regex
          *****************/

        default:
          /* You can pass in /your-own-regex/. */
          if (opts.delimiter instanceof RegExp) {
            delimiterRegex = opts.delimiter;
          } else {
            console.log(NAME + ": Unrecognized delimiter, empty search string, or invalid custom Regex. Aborting.");

            /* Abort this Blast call. */
            return true;
          }
      }
    }

    /**********************
       Element Iteration
    **********************/

    this.each(function() {
      var $this = $(this),
        text = $this.text();

      /* When anything except false is passed in for the options object, Blast is initiated. */
      if (options !== false) {

        /**********************
           Element Variables
        **********************/

        Element = {
          /* The index of each wrapper element generated by blasting. */
          blastedIndex: 0,
          /* Whether we're just entering this node. */
          nodeBeginning: false,
          /* Keep track of the elements generated by Blast so that they can (optionally) be pushed onto the jQuery call stack. */
          wrappers: Element.wrappers || []
        };

        /*****************
           Housekeeping
        *****************/

        /* Unless a consecutive opts.search is being performed, an element's existing Blast call is reversed before proceeding. */
        if ($this.data(NAME) !== undefined && ($this.data(NAME) !== "search" || opts.search === false)) {
          reverse($this, opts);

          if (opts.debug) console.log(NAME + ": Removed element's existing Blast call.");
        }

        /* Store the current delimiter type so that it can be compared against on subsequent calls (see above). */
        $this.data(NAME, opts.search !== false ? "search" : opts.delimiter);

        if (opts.aria) {
          $this.attr("aria-label", text);
        }

        /****************
           Preparation
        ****************/

        /* Perform optional HTML tag stripping. */
        if (opts.stripHTMLTags) {
          $this.html(text);
        }

        /* If the browser throws an error for the provided element type (browers whitelist the letters and types of the elements they accept), fall back to using "span". */
        try {
          document.createElement(opts.tag);
        } catch (error) {
          opts.tag = "span";

          if (opts.debug) console.log(NAME + ": Invalid tag supplied. Defaulting to span.");
        }

        /* For reference purposes when reversing Blast, assign the target element a root class. */
        $this.addClass(NAME + "-root");

        /* Initiate the DOM traversal process. */
        if (opts.debug) console.time(NAME);
        traverseDOM(this, opts);
        if (opts.debug) console.timeEnd(NAME);

        /* If false is passed in as the first parameter, reverse Blast. */
      } else if (options === false && $this.data(NAME) !== undefined) {
        reverse($this, opts);
      }

      /**************
         Debugging
      **************/

      /* Output the full string of each wrapper element and color alternate the wrappers. This is in addition to the performance timing that has already been outputted. */
      if (opts.debug) {
        $.each(Element.wrappers, function(index, element) {
          console.log(NAME + " [" + opts.delimiter + "] " + this.outerHTML);
          this.style.backgroundColor = index % 2 ? "#f12185" : "#075d9a";
        });
      }
    });

    /************
       Reverse
    ************/

    function reverse($this, opts) {
      if (opts.debug) console.time("blast reversal");

      var skippedDescendantRoot = false;

      $this
        .removeClass(NAME + "-root")
        .removeAttr("aria-label")
        .find("." + NAME)
        .each(function() {
          var $this = $(this);
          /* Do not reverse Blast on descendant root elements. (Before you can reverse Blast on an element, you must reverse Blast on any parent elements that have been Blasted.) */
          if (!$this.closest("." + NAME + "-root").length) {
            var thisParentNode = this.parentNode;

            /* This triggers some sort of node layout, thereby solving a node normalization bug in <=IE7 for reasons unknown. If you know the specific reason, tweet me: @Shapiro. */
            if (IE <= 7)(thisParentNode.firstChild.nodeName);

            /* Strip the HTML tags off of the wrapper elements by replacing the elements with their child node's text. */
            thisParentNode.replaceChild(this.firstChild, this);

            /* Normalize() parents to remove empty text nodes and concatenate sibling text nodes. (This cleans up the DOM after our manipulation.) */
            thisParentNode.normalize();
          } else {
            skippedDescendantRoot = true;
          }
        });

      /* Zepto core doesn't include cache-based $.data(), so we mimic data-attr removal by setting it to undefined. */
      if (window.Zepto) {
        $this.data(NAME, undefined);
      } else {
        $this.removeData(NAME);
      }

      if (opts.debug) {
        console.log(NAME + ": Reversed Blast" + ($this.attr("id") ? " on #" + $this.attr("id") + "." : ".") + (skippedDescendantRoot ? " Skipped reversal on the children of one or more descendant root elements." : ""));
        console.timeEnd("blast reversal");
      }
    }

    /*************
        Chain
    *************/

    /* Either return a stack composed of our call's Element.wrappers or return the element(s) originally targeted by the Blast call. */
    /* Note: returnGenerated can only be disabled on a per-call basis (not a per-element basis). */
    if (options !== false && opts.returnGenerated === true) {
      /* A reimplementation of jQuery's $.pushStack() (since Zepto does not provide this function). */
      var newStack = $().add(Element.wrappers);
      newStack.prevObject = this;
      newStack.context = this.context;

      return newStack;
    } else {
      return this;
    }
  };

  /***************
      Defaults
  ***************/

  $.fn.blast.defaults = {
    returnGenerated: true,
    delimiter: "word",
    tag: "span",
    search: false,
    customClass: "",
    generateIndexID: false,
    generateValueClass: false,
    stripHTMLTags: false,
    aria: true,
    debug: false
  };
})(window.jQuery || window.Zepto, window, document);

/*****************
   Known Issues
*****************/

/* In <=IE7, when Blast is called on the same element more than once with opts.stripHTMLTags=false, calls after the first may not target the entirety of the element and/or may
   inject excess spacing between inner text parts due to <=IE7's faulty node normalization. */
