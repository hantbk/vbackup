(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-3acc6b16"],{"3ffe":function(t,e,i){"use strict";i.d(e,"g",(function(){return r})),i.d(e,"e",(function(){return s})),i.d(e,"c",(function(){return n})),i.d(e,"b",(function(){return o})),i.d(e,"q",(function(){return l})),i.d(e,"o",(function(){return c})),i.d(e,"h",(function(){return u})),i.d(e,"n",(function(){return d})),i.d(e,"j",(function(){return p})),i.d(e,"k",(function(){return h})),i.d(e,"a",(function(){return f})),i.d(e,"m",(function(){return m})),i.d(e,"l",(function(){return g})),i.d(e,"i",(function(){return b})),i.d(e,"d",(function(){return y})),i.d(e,"p",(function(){return v})),i.d(e,"f",(function(){return S}));i("99af");var a=i("b775");function r(t){return Object(a["a"])({url:"/repository",method:"get",params:t})}function s(t){return Object(a["a"])({url:"/repository/".concat(t),method:"get"})}function n(t){return Object(a["a"])({url:"/repository/".concat(t),method:"delete"})}function o(t){return Object(a["a"])({url:"/repository",method:"post",data:t})}function l(t){return Object(a["a"])({url:"/repository/".concat(t.id),method:"put",data:t})}function c(t){return Object(a["a"])({url:"/restic/".concat(t.id,"/snapshots"),method:"get",params:t})}function u(t,e,i){return Object(a["a"])({url:"/restic/".concat(t,"/ls/").concat(e),method:"get",params:i})}function d(t,e,i){return Object(a["a"])({url:"/restic/".concat(t,"/search/").concat(e),method:"get",params:i})}function p(t){return Object(a["a"])({url:"/restic/".concat(t,"/parms"),method:"get"})}function h(t){return Object(a["a"])({url:"/restic/".concat(t,"/parmsForMy"),method:"get"})}function f(t){return Object(a["a"])({url:"/restic/".concat(t,"/check"),method:"post"})}function m(t){return Object(a["a"])({url:"/restic/".concat(t,"/rebuild-index"),method:"post"})}function g(t){return Object(a["a"])({url:"/restic/".concat(t,"/prune"),method:"post"})}function b(t){return Object(a["a"])({url:"/restic/".concat(t,"/migrate"),method:"post"})}function y(t,e){return Object(a["a"])({url:"/restic/".concat(t,"/forget"),method:"post",params:{snapshotid:e}})}function v(t,e){return Object(a["a"])({url:"/restic/".concat(t,"/unlock"),method:"post",params:{all:e}})}function S(t,e){return Object(a["a"])({url:"/operation/last/".concat(e,"/").concat(t),method:"get"})}},"56ee":function(t,e,i){"use strict";i("a17a")},a17a:function(t,e,i){},b199:function(t,e,i){"use strict";i.d(e,"c",(function(){return r})),i.d(e,"a",(function(){return s})),i.d(e,"b",(function(){return n}));i("99af");var a=i("b775");function r(t){return Object(a["a"])({url:"/task",method:"get",params:t})}function s(t){return Object(a["a"])({url:"/task/backup/".concat(t),method:"post"})}function n(t,e,i){return Object(a["a"])({url:"/task/".concat(t,"/restore/").concat(e),method:"post",data:i})}},c376:function(t,e,i){"use strict";i.r(e);var a=function(){var t=this,e=t.$createElement,i=t._self._c||e;return i("div",{staticClass:"app-container"},[i("div",{staticClass:"handle-search"},[i("el-form",{attrs:{model:t.listQuery,inline:""},nativeOn:{submit:function(t){t.preventDefault()}}},[i("el-form-item",{attrs:{label:"Filter"}},[i("el-cascader",{staticStyle:{width:"800px"},attrs:{options:t.hostList,props:{expandTrigger:"hover",noDataText:"No data available"},clearable:"",placeholder:"Input content to search",filterable:"",separator:" => "},on:{change:t.handleSearch},model:{value:t.pathQuery,callback:function(e){t.pathQuery=e},expression:"pathQuery"}})],1),i("el-form-item",[i("el-date-picker",{attrs:{type:"date",placeholder:"Select date"},on:{change:t.handleSearch},model:{value:t.listQuery.date,callback:function(e){t.$set(t.listQuery,"date",e)},expression:"listQuery.date"}})],1)],1)],1),i("div",[i("el-row",{staticClass:"panel-group",attrs:{gutter:40}},[i("el-col",{staticClass:"card-panel-col",attrs:{xs:8,sm:8,lg:8}},[i("el-card",{staticClass:"box-card"},[i("div",{attrs:{slot:"header"},slot:"header"},[i("p",[t._v("Host: "+t._s(t.listQuery.host+(t.list.length>0?"("+t.list[0].username+")":"")))]),i("p",[t._v("Path: "+t._s(t.listQuery.path))]),i("p",[t._v("Total: "+t._s(t.total))])]),i("el-collapse",{attrs:{accordion:""},model:{value:t.activeName,callback:function(e){t.activeName=e},expression:"activeName"}},t._l(t.snapList,(function(e,a){return i("el-collapse-item",{key:a,attrs:{title:e.name,name:e.name}},[i("el-timeline",t._l(e.list,(function(e,a){return i("el-timeline-item",{key:a,attrs:{timestamp:t._f("goDatToDateString")(e.time),type:"primary",icon:"el-icon-success",size:"large",placement:"top"}},[i("div",{staticClass:"timeline"},[i("span",{staticClass:"snap",on:{click:function(i){return t.loadSnapFiles(e)}}},[t._v(t._s(e.short_id))]),i("el-button",{staticClass:"restore_btn",attrs:{type:"text",size:"mini"},on:{click:function(i){return t.openRestoreOpt(e.short_id)}}},[t._v(" Restore ")])],1)])})),1)],1)})),1),t.noMore?i("p",{staticStyle:{"font-size":"20px","text-align":"center",color:"#bbbbbb"}},[t._v("No more")]):i("div",{staticStyle:{"text-align":"center"}},[i("el-button",{attrs:{loading:t.listLoading,type:"info",plain:""},on:{click:t.getMoreList}},[t._v("Load more")])],1)],1)],1),i("el-col",{staticClass:"card-panel-col",attrs:{xs:16,sm:16,lg:16}},[i("el-card",{staticClass:"box-card"},[i("div",[i("el-input",{staticClass:"input-with-select",attrs:{placeholder:"Enter the exact path for faster search, e.g: data/test/avatar.png",clearable:""},on:{clear:t.searchFile},model:{value:t.fileSearch.name,callback:function(e){t.$set(t.fileSearch,"name",e)},expression:"fileSearch.name"}},[i("el-select",{staticStyle:{width:"170px"},attrs:{slot:"prepend",placeholder:"Please select"},slot:"prepend",model:{value:t.fileSearch.type,callback:function(e){t.$set(t.fileSearch,"type",e)},expression:"fileSearch.type"}},t._l(t.searchType,(function(t,e){return i("el-option",{key:e,attrs:{label:t.name,value:t.code}})})),1),i("el-button",{attrs:{slot:"append",icon:"el-icon-search"},on:{click:t.searchFile},slot:"append"})],1)],1),i("el-tree",{directives:[{name:"loading",rawName:"v-loading",value:t.treeLoading,expression:"treeLoading"}],staticClass:"file-tree",attrs:{data:t.filedata.children,"node-key":"id",accordion:"","empty-text":"No data available","expand-on-click-node":""},on:{"node-expand":t.getFiles},scopedSlots:t._u([{key:"default",fn:function(e){var a=e.node,r=e.data;return i("span",{staticClass:"custom-tree-node",on:{click:function(e){return t.moreClick(r,a)}}},[i("div",{staticClass:"file-title"},["dir"===r.type?i("i",{staticClass:"el-icon-folder"}):t._e(),"btn"===r.type?i("i",{staticClass:"el-icon-more-outline"}):t._e(),"file"===r.type?i("i",{staticClass:"el-icon-document"}):t._e(),i("span",{staticStyle:{"margin-left":"5px"}},[t._v(t._s(a.label))])]),i("span",[i("span",{staticStyle:{"margin-right":"10px","font-size":"13px",display:"inline"}},[t._v(t._s(r.size))]),r.loading?i("i",{staticClass:"el-icon-loading",attrs:{type:"primary"}}):t._e(),0!==r.isMore||r.loading?t._e():i("el-button",{attrs:{type:"text",size:"mini"},on:{click:function(){return t.restoreFileHandler(r)}}},[t._v(" Restore ")])],1)])}}])})],1)],1)],1)],1),i("el-dialog",{attrs:{title:"Restore Options",visible:t.dialogFormVisible},on:{"update:visible":function(e){t.dialogFormVisible=e}}},[i("el-form",{ref:"dataForm",attrs:{"label-position":"left","label-width":"260px"}},[i("el-form-item",{attrs:{label:"Restore to: ",prop:"path"}},[i("el-input",{attrs:{disabled:""},model:{value:t.restoreOpt.dirCur,callback:function(e){t.$set(t.restoreOpt,"dirCur",e)},expression:"restoreOpt.dirCur"}},[i("el-button",{attrs:{slot:"append"},on:{click:function(e){return t.openDirSelect()}},slot:"append"},[t._v("Select")])],1),i("span",{staticStyle:{color:"red"}},[t._v("Default restore to '/', restoring data to the original path of the file. If modified, the data restore path will be the current selected path plus the backup path, for example: /root"+t._s(t.listQuery.path)+", /root is the current selected path")])],1),i("el-form-item",{attrs:{label:"Final directory for the data: ",prop:"path"}},[i("span",[t._v(t._s(t.restoreOpt.dist))])]),i("el-form-item",{attrs:{label:"Do you want to verify data integrity?"}},[i("el-switch",{model:{value:t.restoreOpt.verify,callback:function(e){t.$set(t.restoreOpt,"verify",e)},expression:"restoreOpt.verify"}}),i("p",{staticStyle:{color:"red"}},[t._v("Enable data integrity verification, which may take a long time. Please choose according to your needs! ")])],1)],1),i("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[i("el-button",{on:{click:function(e){t.dialogFormVisible=!1}}},[t._v(" Cancel ")]),i("el-button",{attrs:{type:"primary"},on:{click:function(e){return t.restoreSnapHandler()}}},[t._v(" Confirm ")])],1)],1),i("el-dialog",{attrs:{title:"Select Folder",visible:t.dialogDirVisible},on:{"update:visible":function(e){t.dialogDirVisible=e}}},[i("div",[i("el-breadcrumb",{attrs:{separator:"/"}},t._l(t.getDirSpea(),(function(e,a){return i("el-breadcrumb-item",{key:a,staticClass:"breadcrumb-item"},[i("span",{staticClass:"title",on:{click:function(i){return t.lsDir(e.path,!0)}}},[t._v(t._s(e.name))])])})),1),i("div",{staticClass:"filenodes"},t._l(t.filteredDirList,(function(e,a){return i("span",{key:a,staticClass:"custom-tree-node filenode",class:{active:t.restoreOpt.dirCur===e.path},on:{dblclick:function(i){return i.preventDefault(),t.lsDir(e.path,e.isDir)},click:function(i){return t.selectDir(e.path,e.isDir)}}},[i("div",{staticClass:"file-title"},[e.isDir?i("i",{staticClass:"el-icon-folder"}):t._e(),i("span",{staticStyle:{"margin-left":"5px","user-select":"none"}},[t._v(t._s(e.name))])]),i("span",[i("el-button",{staticClass:"confirmbtn",attrs:{type:"text",size:"mini"},on:{click:function(i){return t.confirmDirSelect(e.path)}}},[t._v(" Confirm ")])],1)])})),0)],1),i("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[i("el-button",{on:{click:function(e){t.dialogDirVisible=!1}}},[t._v(" Cancel ")]),i("el-button",{attrs:{type:"primary"},on:{click:function(e){return t.confirmDirSelect()}}},[t._v(" Confirm ")])],1)])],1)},r=[],s=(i("99af"),i("4de4"),i("b0c0"),i("a9e3"),i("d3b7"),i("159b"),i("3ffe")),n=i("ed08"),o=i("b199"),l=i("8593"),c={name:"Restore",data:function(){return{list:[],total:0,snapList:[],tagList:[],hostList:[],dialogFormVisible:!1,dialogDirVisible:!1,dirList:[],restoreOpt:{dirCur:"/",snapid:"",include:"",dist:"",verify:!1},listLoading:!1,treeLoading:!1,noMore:!0,curSnap:{},activeName:"0",searchType:[{code:0,name:"Full Library"},{code:1,name:"Current Snapshot"}],fileSearch:{type:0,name:"",pageNum:1,pageSize:20},pathQuery:[],listQuery:{id:0,path:"",date:"",host:"",tags:"",pageNum:1,pageSize:100},filedata:{label:"root",path:"",type:"dir",isMore:3,pageNum:1,pageSize:20,children:[]}}},created:function(){this.listQuery.id=this.$route.params&&this.$route.params.id,this.getParmList()},activated:function(){this.getParmList()},computed:{filteredDirList:function(){return this.dirList.filter((function(t){return t.isDir}))}},methods:{openRestoreOpt:function(t){this.dialogFormVisible=!0,this.restoreOpt.snapid=t,this.restoreOpt.dirCur="/",this.restoreOpt.dist=this.listQuery.path,this.restoreOpt.include=this.listQuery.path},restoreFileHandler:function(t){this.dialogFormVisible=!0,this.restoreOpt.snapid=this.curSnap.short_id,this.restoreOpt.dirCur="/",this.restoreOpt.dist=t.path,this.restoreOpt.include=t.path},openDirSelect:function(){this.dialogDirVisible=!0,this.dirList=[],this.lsDir(this.restoreOpt.dirCur,!0)},confirmDirSelect:function(t){t&&(this.restoreOpt.dirCur=t,this.restoreOpt.dist=t+this.restoreOpt.include),this.dialogDirVisible=!1},getDirSpea:function(){var t=this.restoreOpt.dirCur.split("/");t.shift();var e=[],i="";return t.forEach((function(t){""!==t&&(i=i+"/"+t,e.push({name:t,path:i}))})),e.unshift({name:"Root",path:"/"}),e},selectDir:function(t,e){e&&(this.restoreOpt.dirCur=t,this.restoreOpt.dist=t+this.restoreOpt.include)},lsDir:function(t,e){var i=this;if(e){this.restoreOpt.dirCur=t,this.restoreOpt.dist=t+this.restoreOpt.include;var a={path:this.restoreOpt.dirCur};Object(l["b"])(a).then((function(t){i.dirList=t.data}))}},getMoreList:function(){this.listQuery.pageNum++,this.getList()},handleSearch:function(){this.listQuery.host="",this.listQuery.path="",this.listQuery.pageNum=1,this.list=[],this.filedata={label:"root",path:"",type:"dir",isMore:3,pageNum:1,pageSize:20,children:[]},this.fileSearch={type:0,name:"",pageNum:1,pageSize:20},this.snapList=[],2===this.pathQuery.length&&(this.listQuery.host=this.pathQuery[0],this.listQuery.path=this.pathQuery[1],this.getList())},searchFile:function(){var t=this;if(1===this.fileSearch.type){if(""===this.fileSearch.name)return void this.loadSnapFiles(this.curSnap);var e={path:this.fileSearch.name,pageNum:this.fileSearch.pageNum,pageSize:this.fileSearch.pageSize};this.filedata.children=[],this.treeLoading=!0,Object(s["n"])(this.listQuery.id,this.curSnap.short_id,e).then((function(e){var i=e.data.pageNum,a=e.data.pageSize,r=e.data.items.nodes;r.forEach((function(e){var r={pageNum:i,pageSize:a,name:e.name,path:e.path,label:e.path,type:e.type,mode:e.mode,isMore:0,permissions:e.permissions,ctime:e.ctime,gid:e.gid,uid:e.uid,size:e.size,children:[]};t.filedata.children.push(r)}))})).finally((function(){t.treeLoading=!1}))}else this.$notify.error("This feature is not yet available.")},restoreSnapHandler:function(){var t=this,e=this.restoreOpt.snapid;this.$confirm("Are you sure you want to execute the restore operation for <"+e+">? This operation may take a long time! ","Restore Data",{type:"warning"}).then((function(){var i={target:t.restoreOpt.dirCur,include:t.restoreOpt.include,verify:t.restoreOpt.verify};Object(o["b"])(t.listQuery.id,e,i).then((function(e){t.dialogDirVisible=!1,t.dialogFormVisible=!1,t.$notify.success({title:"Restoring...",message:'Please go to "<a style="color: #409EFF" href="/Task/index">Task Records</a>" to check.',dangerouslyUseHTMLString:!0})})).finally((function(){t.dialogDirVisible=!1,t.dialogFormVisible=!1}))})).catch((function(){t.$notify.info({title:"Cancel"})}))},loadSnapFiles:function(t){this.fileSearch.type=1,this.curSnap=t,this.filedata.pageNum=1,this.filedata.isMore=3,this.filedata.path=this.listQuery.path,this.filedata.children=[{pageNum:this.filedata.pageNum,pageSize:this.filedata.pageSize,name:"",path:this.filedata.path,label:"Load more",isMore:1,type:"btn",mode:755,permissions:"",ctime:"",gid:"",uid:"",size:"",children:[]}],this.getFiles(this.filedata,{parent:{data:{children:[],pageNum:1}}})},moreClick:function(t,e){"btn"===t.type&&this.getFiles(t,e)},getFiles:function(t,e){var i=this;if((0!==t.isMore||"file"!==t.type)&&2!==t.isMore){"dir"===t.type?t.pageNum=1:t.pageNum++;var a={path:t.path,pageNum:t.pageNum,pageSize:t.pageSize};this.treeLoading=!0,Object(s["h"])(this.listQuery.id,this.curSnap.short_id,a).then((function(a){var r,s=a.data.pageNum,n=a.data.pageSize,o=a.data.total,l=a.data.items.nodes;if(l.forEach((function(a){var r=[];"dir"===a.type&&(r=[{pageNum:1,pageSize:n,path:a.path,name:"",label:"Load more",isMore:1,type:"btn",mode:a.mode,permissions:a.permissions,ctime:a.ctime,gid:a.gid,uid:a.uid,size:a.size,children:[]}]);var s={pageNum:1,pageSize:n,name:a.name,path:a.path,label:a.name,isMore:0,type:a.type,mode:a.mode,permissions:a.permissions,ctime:a.ctime,gid:a.gid,uid:a.uid,size:a.size,children:r};if(1===t.isMore){var o=e.parent;o.data.children?o.data.children.push(s):i.filedata.children.push(s)}else t.children.push(s)})),"file"!==t.type)if(1===t.isMore){var c=e.parent;c.data.children?(r=c.data.children,c.data.children=i.moveMoretoLast(r,s,n,o),Number(s)*Number(n)>=o&&(c.data.isMore=0)):(r=i.filedata.children,i.filedata.children=i.moveMoretoLast(r,s,n,o))}else r=t.children,t.children=i.moveMoretoLast(r,s,n,o),Number(s)*Number(n)>=o&&(t.isMore=0)})).finally((function(){i.treeLoading=!1}))}},moveMoretoLast:function(t,e,i,a){var r=[],s=null;return t.forEach((function(t){if(1===t.isMore)return s=t,void(Number(e)*Number(i)>=a&&(s.isMore=2,s.label="No more"));r.push(t)})),r.push(s),r},getParmList:function(){var t=this;Object(s["j"])(this.listQuery.id).then((function(e){t.tagList=e.data.tags,t.hostList=t.buildCascader(e.data.parms)}))},buildCascader:function(t){var e=[];return t.forEach((function(t){var i=[];t.paths.forEach((function(t){var e={value:t,label:t};i.push(e)}));var a={value:t.name,label:t.name,children:i};e.push(a)})),e},getList:function(){var t=this;this.listLoading=!0,Object(s["o"])(this.listQuery).then((function(e){t.list=t.list.concat(e.data.items),t.snapList=t.accordionList(t.list),t.noMore=Number(e.data.pageNum)*Number(e.data.pageSize)>=e.data.total,t.total=e.data.total})).finally((function(){t.listLoading=!1}))},accordionList:function(t){var e=this,i=[];return t.forEach((function(t){var a=Object(n["a"])(t.time,"yyyy-MM");if(e.listQuery.date){var r=Object(n["a"])(t.time,"yyyy-MM-dd"),s=Object(n["a"])(e.listQuery.date,"yyyy-MM-dd");if(r!==s)return;e.activeName=a}var o=!1;i.forEach((function(e){e.name===a&&(e.list.push(t),o=!0)})),o||i.push({name:a,list:[t]})})),i}}},u=c,d=(i("56ee"),i("2877")),p=Object(d["a"])(u,a,r,!1,null,"36f4d1b6",null);e["default"]=p.exports}}]);