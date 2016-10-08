/* Copyright 2015 The TensorFlow Authors. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
==============================================================================*/

/**
 * @fileoverview This file manually reloads the main content in templated pages,
 * so that any state (e.g. how the sitemap accordion is open) will persist
 */

window.onload = function() {
  if (!window.XMLHttpRequest) {
    return;  // for really really old browsers. nojs site will still work.
  }
  var sitemap = document.getElementById('sitemap');
  var content = document.getElementById('content');

  function getBase() {
    // Only create the base when actually needed (since I don't want to
    // figure out what the right relative path is before this script has done
    // anything, and base without a target or href is invalid)
    var base = document.getElementsByTagName('base')[0];
    if (!base) {
      base = document.createElement('base');
      document.getElementsByTagName('head')[0].appendChild(base);
    }
    return base;
  }

  function reloadContent(path, body_path, hash){
    return function() {
      var httpRequest = new XMLHttpRequest();
      httpRequest.onreadystatechange = function() {
        if (httpRequest.readyState === XMLHttpRequest.DONE) {
          getBase().href = path;  // So relative paths load correctly
          content.innerHTML = httpRequest.responseText;
          // Get the new title and set it
          headers = content.getElementsByTagName('H1');
          if (headers.length > 0) {
            window.document.title = 'TensorFlow -- ' + headers[0].textContent;
          }

          function navigateToAnchor() {
            // Defer changing hash and navigating to anchor until after
            // MathJax is finished loading, so there are no surprising offsets
            // from top.
            location.hash = hash;
            // Correct history entry.
            window.history.replaceState(
                window.history.state, window.document.title, path + hash);
          }
          MathJax.Hub.Queue(['Typeset', MathJax.Hub]);  // Update MathJax
          MathJax.Hub.Queue(navigateToAnchor);
        }
      };
      httpRequest.open('GET', body_path, true);
      httpRequest.send(null);
      return false;
    }
  };
  var as = sitemap.getElementsByTagName('a');
  for (var i = 0; i < as.length; i++) {
    var a = as[i];
    var path = a.pathname;
    body_path = path.replace('.html', '_body.html');
    a.onclick = reloadContent(path, body_path, a.hash);
  }
}
