const pug = require('pug')

var tt = function(filename, list, indentLevel = 0, hrefParent = '/documentation') {
  var str = '';
  var indentStr = '  ';
  var indentation = indentStr.repeat(indentLevel);
  for (var i = 0, len = list.length; i < len; i++) {
    var href = hrefParent + (list[i].href != undefined ? '/' + list[i].href : '')
    switch (list[i].type) {
      case 'label':
        str += indentation + 'p.menu-label ' + list[i].text + '\r\n';
        break;
      case 'list':
        str += indentation + 'ul.menu-list' + '\r\n';
        break;
      case 'item':
        var anchor = list[i].anchor != undefined ? '#' + list[i].anchor : ''
        var cssClass = 'tta-doc-menu-item'
                        + (filename.replace('.pug', '.html').substr(1) == href + '.html' + anchor ? '.is-active' : '')
        str += indentation
                + 'li\r\n'
                + indentation + indentStr + 'a.' + cssClass + '(href="' + href + '.html' + anchor + '") ' + list[i].text
                + '\r\n';
        break;
    }
    if (list[i].itemList != undefined && list[i].itemList.length > 0) {
      str += tt(filename, list[i].itemList, indentLevel + 1, href);
    }
  }
  return str;
}

var tt1 = function(filename, list, indentLevel = 0) {
  return pug.render(tt(filename, list, indentLevel));
}

exports.indenter = tt1
