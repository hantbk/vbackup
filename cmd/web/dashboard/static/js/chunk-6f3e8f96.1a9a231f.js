(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-6f3e8f96"],{"2cbf":function(e,t,a){"use strict";a("73e0")},"333d":function(e,t,a){"use strict";var n=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"pagination-container",class:{hidden:e.hidden}},[a("el-pagination",e._b({attrs:{background:e.background,"current-page":e.currentPage,"page-size":e.pageSize,layout:e.layout,"page-sizes":e.pageSizes,total:e.total},on:{"update:currentPage":function(t){e.currentPage=t},"update:current-page":function(t){e.currentPage=t},"update:pageSize":function(t){e.pageSize=t},"update:page-size":function(t){e.pageSize=t},"size-change":e.handleSizeChange,"current-change":e.handleCurrentChange}},"el-pagination",e.$attrs,!1))],1)},r=[];a("a9e3");Math.easeInOutQuad=function(e,t,a,n){return e/=n/2,e<1?a/2*e*e+t:(e--,-a/2*(e*(e-2)-1)+t)};var i=function(){return window.requestAnimationFrame||window.webkitRequestAnimationFrame||window.mozRequestAnimationFrame||function(e){window.setTimeout(e,1e3/60)}}();function l(e){document.documentElement.scrollTop=e,document.body.parentNode.scrollTop=e,document.body.scrollTop=e}function o(){return document.documentElement.scrollTop||document.body.parentNode.scrollTop||document.body.scrollTop}function u(e,t,a){var n=o(),r=e-n,u=20,s=0;t="undefined"===typeof t?500:t;var c=function(){s+=u;var e=Math.easeInOutQuad(s,n,r,t);l(e),s<t?i(c):a&&"function"===typeof a&&a()};c()}var s={name:"Pagination",props:{total:{required:!0,type:Number},page:{type:Number,default:1},limit:{type:Number,default:20},pageSizes:{type:Array,default:function(){return[10,20,30,50]}},layout:{type:String,default:"total, sizes, prev, pager, next, jumper"},background:{type:Boolean,default:!0},autoScroll:{type:Boolean,default:!0},hidden:{type:Boolean,default:!1}},computed:{currentPage:{get:function(){return this.page},set:function(e){this.$emit("update:page",e)}},pageSize:{get:function(){return this.limit},set:function(e){this.$emit("update:limit",e)}}},methods:{handleSizeChange:function(e){this.$emit("pagination",{page:this.currentPage,limit:e}),this.autoScroll&&u(0,800)},handleCurrentChange:function(e){this.$emit("pagination",{page:e,limit:this.pageSize}),this.autoScroll&&u(0,800)}}},c=s,d=(a("2cbf"),a("2877")),p=Object(d["a"])(c,n,r,!1,null,"6af373ef",null);t["a"]=p.exports},5905:function(e,t,a){"use strict";a.r(t);var n=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"app-container"},[a("div",{staticClass:"handle-search"},[a("el-form",{attrs:{model:e.listQuery,inline:""},nativeOn:{submit:function(e){e.preventDefault()}}},[a("el-form-item",{attrs:{label:"Operator"}},[a("el-input",{staticClass:"filter-item",attrs:{placeholder:"Operator",clearable:""},model:{value:e.listQuery.operator,callback:function(t){e.$set(e.listQuery,"operator",t)},expression:"listQuery.operator"}})],1),a("el-form-item",{attrs:{label:"Action"}},[a("el-select",{staticClass:"handle-select mr5",attrs:{clearable:"",placeholder:"Please select"},model:{value:e.listQuery.operation,callback:function(t){e.$set(e.listQuery,"operation",t)},expression:"listQuery.operation"}},e._l([{value:"",name:"All"}].concat(e.operationList),(function(e,t){return a("el-option",{key:t,attrs:{label:e.name,value:e.value}})})),1)],1),a("el-form-item",{attrs:{label:"Resource"}},[a("el-select",{staticClass:"handle-select mr5",attrs:{clearable:"",placeholder:"Please select"},model:{value:e.listQuery.url,callback:function(t){e.$set(e.listQuery,"url",t)},expression:"listQuery.url"}},e._l([{value:"",name:"All"}].concat(e.resourceList),(function(e,t){return a("el-option",{key:t,attrs:{label:e.name,value:e.value}})})),1)],1),a("el-form-item",{attrs:{label:"Data"}},[a("el-input",{staticClass:"filter-item",attrs:{placeholder:"",clearable:""},model:{value:e.listQuery.data,callback:function(t){e.$set(e.listQuery,"data",t)},expression:"listQuery.data"}})],1),a("el-form-item",[a("el-button",{attrs:{type:"primary",icon:"el-icon-search"},on:{click:e.handleFilter}},[e._v("Search")])],1)],1)],1),a("el-table",{directives:[{name:"loading",rawName:"v-loading",value:e.listLoading,expression:"listLoading"}],staticStyle:{width:"100%"},attrs:{data:e.list,border:"",fit:"","highlight-current-row":""}},[a("el-table-column",{attrs:{prop:"id",align:"center",label:"ID"}}),a("el-table-column",{attrs:{prop:"createdAt",formatter:e.dateFormat,align:"center",label:e._f("i18n")("createdAt")}}),a("el-table-column",{attrs:{prop:"operator",align:"center",label:"Operator"}}),a("el-table-column",{attrs:{prop:"operation",align:"center",label:"Action"},scopedSlots:e._u([{key:"default",fn:function(t){var n=t.row;return[a("el-tag",{attrs:{type:e.filterOper(n.operation).color}},[e._v(" "+e._s(e.filterOper(n.operation).name)+" ")])]}}])}),a("el-table-column",{attrs:{prop:"url",align:"center",label:"Resource"},scopedSlots:e._u([{key:"default",fn:function(t){var a=t.row;return[e._v(" "+e._s(e.filterResource(a.url).name)+" ")]}}])}),a("el-table-column",{attrs:{prop:"url",align:"center",label:"Address"}}),a("el-table-column",{attrs:{prop:"data",align:"center",label:"Data"},scopedSlots:e._u([{key:"default",fn:function(t){var n=t.row;return[n.data?a("i",{staticClass:"el-icon-document",on:{click:function(t){return e.handleInfo(n)}}}):e._e()]}}])})],1),a("pagination",{directives:[{name:"show",rawName:"v-show",value:e.total>0,expression:"total > 0"}],attrs:{total:e.total,page:e.listQuery.pageNum,limit:e.listQuery.pageSize,autoScroll:!1},on:{"update:page":function(t){return e.$set(e.listQuery,"pageNum",t)},"update:limit":function(t){return e.$set(e.listQuery,"pageSize",t)},pagination:e.getList}}),a("el-dialog",{attrs:{title:"Data Details",visible:e.dialogFormVisible,top:"5vh"},on:{"update:visible":function(t){e.dialogFormVisible=t}}},[a("el-input",{staticClass:"json-text",attrs:{type:"textarea",readonly:"",autosize:""},model:{value:e.flowJSON,callback:function(t){e.flowJSON=t},expression:"flowJSON"}})],1)],1)},r=[],i=(a("7db0"),a("caad"),a("e9c4"),a("b64b"),a("d3b7"),a("2532"),a("5fd4")),l=a("ed08"),o=a("333d"),u={name:"LogList",components:{Pagination:o["a"]},data:function(){return{operationList:[{name:"Add",value:"post",color:"success"},{name:"Edit",value:"put",color:"primary"},{name:"Delete",value:"delete",color:"danger"}],resourceList:[{name:"Login",value:"login"},{name:"Repository",value:"repository"},{name:"Backup Plan",value:"plan"},{name:"User",value:"user"},{name:"Execute Backup",value:"backup"},{name:"Execute Data Restore",value:"restore"},{name:"Cleanup Policy",value:"policy"},{name:"Data Maintenance",value:"restic"}],listLoading:!1,listQuery:{operator:"",operation:"",url:"",data:"",pageNum:1,pageSize:10},list:[],total:0,dialogFormVisible:!1,flowJSON:{}}},created:function(){this.getList()},methods:{dateFormat:function(e,t,a,n){return Object(l["a"])(a,"yyyy-MM-dd hh:mm:ss")},filterOper:function(e){return this.operationList.find((function(t){return t.value===e}))||{name:"",value:e}},filterResource:function(e){return this.resourceList.find((function(t){return e.includes(t.value)}))||{name:"",value:e}},handleFilter:function(){this.listQuery.pageNum=1,this.getList()},handleInfo:function(e){this.dialogFormVisible=!0,this.flowJSON={},this.flowJSON=JSON.stringify(JSON.parse(e.data),null,"    ")},getList:function(){var e=this;this.listLoading=!0,Object(i["c"])(this.listQuery).then((function(t){e.list=t.data.items,e.total=t.data.total})).finally((function(){e.listLoading=!1}))}}},s=u,c=(a("a215"),a("2877")),d=Object(c["a"])(s,n,r,!1,null,"0dbe4e36",null);t["default"]=d.exports},"5fd4":function(e,t,a){"use strict";a.d(t,"b",(function(){return r})),a.d(t,"a",(function(){return i})),a.d(t,"c",(function(){return l}));var n=a("b775");function r(e){return Object(n["a"])({url:"/dashboard/index",method:"get",params:e})}function i(){return Object(n["a"])({url:"/dashboard/doGetAllRepoStats",method:"post"})}function l(e){return Object(n["a"])({url:"/dashboard/logs",method:"get",params:e})}},"73e0":function(e,t,a){},"7db0":function(e,t,a){"use strict";var n=a("23e7"),r=a("b727").find,i=a("44d2"),l=a("ae40"),o="find",u=!0,s=l(o);o in[]&&Array(1)[o]((function(){u=!1})),n({target:"Array",proto:!0,forced:u||!s},{find:function(e){return r(this,e,arguments.length>1?arguments[1]:void 0)}}),i(o)},a215:function(e,t,a){"use strict";a("e3af")},e3af:function(e,t,a){},e9c4:function(e,t,a){var n=a("23e7"),r=a("d066"),i=a("d039"),l=r("JSON","stringify"),o=/[\uD800-\uDFFF]/g,u=/^[\uD800-\uDBFF]$/,s=/^[\uDC00-\uDFFF]$/,c=function(e,t,a){var n=a.charAt(t-1),r=a.charAt(t+1);return u.test(e)&&!s.test(r)||s.test(e)&&!u.test(n)?"\\u"+e.charCodeAt(0).toString(16):e},d=i((function(){return'"\\udf06\\ud834"'!==l("\udf06\ud834")||'"\\udead"'!==l("\udead")}));l&&n({target:"JSON",stat:!0,forced:d},{stringify:function(e,t,a){var n=l.apply(null,arguments);return"string"==typeof n?n.replace(o,c):n}})}}]);