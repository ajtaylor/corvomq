pug = require('pug')

var tt = function(list, indentLevel = 0, hrefParent = '\documentation') {
  var str = '';
  var indentStr = '  ';
  var indentation = indentStr.repeat(indentLevel);
  for (var i = 0, len = list.length; i < len; i++) {
    switch (list[i].type) {
      case 'label':
        str += indentation + 'p.menu-label ' + list[i].text + '\r\n';
        break;
      case 'list':
        str += indentation + 'ul.menu-list' + '\r\n';
        break;
      case 'item':
        str += indentation
                + 'li\r\n'
                + indentation + indentStr + 'a.tta-doc-menu-item ' + list[i].text
                + '\r\n';
        break;
    }
    if (list[i].itemList != undefined && list[i].itemList.length > 0) {
      str += tt(list[i].itemList, indentLevel + 1);
    }
  }
  return str;
}

var tt1 = function(list, indentLevel = 0) {
  return pug.render(tt(list, indentLevel));
}

exports.indenter = tt1
