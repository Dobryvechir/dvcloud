package dvmongodb

import (
	"strings"
)

const REF_FIELD = "_id"
const ID_FIELD = "id"
const WHERE_START = "{$where:\""
const WHERE_FINISH = "\"}"

var MongoFunctions map[string][]string = map[string][]string{
	"round": []string{"Math.round(", ")"},
	"ceil":  []string{"Math.ceil(", ")"},
	"floor": []string{"Math.floor(", ")"},
	"min":   []string{"Math.min(", ")"},
	"max":   []string{"Math.max(", ")"},
	"year":  []string{"dd(", ").getFullYear()"},
	"month": []string{"(dd(", ").getMonth()+1)"},
	"day":   []string{"dd(", ").getDate()"}}

const READ_FIELD_IN_ARRAY = WHERE_START + "function(){var t=this,a=function(l){var n=l.length,b=[t],r,e,i,g,h;" +
	"for(i=0;i<n;i++){" +
	"e=l[i];h=[];b.forEach(function(d){" +
	"g=d[e];if(!g && e=='$id' && g.getId)g=d.getId();" +
	"if (g!==undefined && g!==null){" +
	"if (Array.isArray(g)){g.forEach(function(f){" +
	"if (f!==undefined && f!==null)h.push(f);})}" +
	"else h.push(g);}});" +
	"b=h;if(!b.length)return null;" +
	"}return b;},"
const PROCESS_ARITHMETIC_EXPRESSION = READ_FIELD_IN_ARRAY +
	"s=function(x){return typeof x==='string'?x:x+'';}," +
	"d=function(x){return new Date(x).getTime()/86400000;},dd=function(x){return new Date(Math.round(x*86400000));}," +
	"v=function(x){if(x===null||x===undefined)return null;if(typeof x!=='object')return x;" +
	"if(x.getId)return x.getId();if(x.getTime)return x.getTime()/86400000;var y=x.toString(),z=y.indexOf('(');" +
	"if(z>0&&y.startsWith('Number')){if(y.charCodeAt(z+1)===34)return +y.substring(z+2,y.length-2);return x.valueOf();}return y;}," +
	"b=function(x){var i,n=x.length,j=[],m=[],y=[],p=-1,r;for(i=0;i<n;i++){if(!x[i]||!x[i].length)return false;j[i]=0;m[i]=x[i].length;}" +
	"while(p<n){p=0;r=0;" +
	"for(i=0;i<n;i++){" +
	"y[i]=v(x[i][j[i]]);if(y[i]===null)r=1;" +
	"if (i===p&&(++j[i])>=m[i]){p++;j[i]=0;}}" +
	"if (r==0&&("
const PROCESS_ARITHMETIC_EXPRESSION_MIDDLE = "))return true;}return false;};return b(["
const PROCESS_ARITHMETIC_EXPRESSION_END = "]);}" + WHERE_FINISH

func createGeneralWhereExpression(fields []string, expression string) string {
	result := PROCESS_ARITHMETIC_EXPRESSION + expression + PROCESS_ARITHMETIC_EXPRESSION_MIDDLE
	for k, v := range fields {
		if k != 0 {
			result += ","
		}
		result += "a([" + appendField(v) + "])"
	}
	return result + PROCESS_ARITHMETIC_EXPRESSION_END
}

func createWhereForInOperations(field string, fields []string, isReverse bool) string {
	result := READ_FIELD_IN_ARRAY + "p=a([" + appendField(field) + "]),s=["
	for k, v := range fields {
		if k != 0 {
			result += ","
		}
		result += "a([" + v + "])"
	}
	reverse := "!!"
	if isReverse {
		reverse = "!"
	}
	return result + "],b={},m=p&&p.length,i;if(!m)return" + reverse +
		"0;s.forEach(function(x){if (x) x.forEach(function(y){b[y]=1;})});" +
		"for(i=0;i<m;i++)if(b[p[i]]) return" + reverse +
		"1;return" + reverse + "0;}" + WHERE_FINISH
}

func appendField(field string) (res string) {
	r := strings.Split(field, ".")
	res = ""
	for k, v := range r {
		if k != 0 {
			res += ","
		}
		res += "'" + v + "'"
	}
	return
}
