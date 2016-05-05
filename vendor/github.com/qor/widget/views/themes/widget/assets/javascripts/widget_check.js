if(!window.loadedWidgetAsset) {
  window.loadedWidgetAsset = true;
  var prefix = document.currentScript.getAttribute("data-prefix");
  document.write("<script src='" + prefix + "/assets/javascripts/vendors/jquery.min.js'></script><script src=\"" + prefix + "/assets/javascripts/widget.js?theme=widget\"></script><link type=\"text/css\" rel=\"stylesheet\" href=\"" + prefix + "/assets/stylesheets/widget.css?theme=widget\">");
}
