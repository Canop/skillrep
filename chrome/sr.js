
// things

function commated(n){
	return (''+n).replace(/(^\d{1,3}|\d{3})(?=(?:\d{3})+(?:$|\.))/g, '$1,')
}

function enrichUserInfo(parentNode){
	parentNode.querySelectorAll(".post-signature").forEach(function(ui){
		console.log("user-info:", ui.innerHTML)
		var userLink = ui.querySelector(".user-info .user-details a")
		if (!userLink) return
		var userId = userLink.href.match(/\/users\/(\d+)/)[1]
		console.log("UID:",userId)
		fetch("https://dystroy.org/skillrep/users/"+userId).then(function(data){
			console.log("received:", data)
			ui.style.border = "2px solid #444"
			var srbox = ui.append('a')
			srbox.className = 'skillrep-user'
			srbox.href = 'https://dystroy.org/skillrep/?user='+userId
			srbox.target = "_blank"
			srbox.append('<div class=skillrep-title>', "SkillRep")
			if (data.User.Rank) {
				srbox.append('<div class=skillrep-rank>', commated(data.User.Rank))
			}
			srbox.append('<div class=skillrep-rep>', commated(data.User.SkillRep||0))
		})
	})

}
var list = document.querySelector('questions')
if (list) {
	(new MutationObserver(function(mutations){
		for (var i=0; i<mutations.length; i++) {
			var nodes = mutations[i].addedNodes
			if (!nodes) continue
			nodes.forEach(enrichUserInfo)
		}
	})).observe(document.getElementById('questions'), {childList:true})
}

function fetch(path,args){
	if (args) {
		path += (/\?/.test(path) ? '&' : '?') + Object.keys(args).map(function(k){
			return k + '=' + encodeURIComponent(args[k])
		}).join('&')
	}
	return new Promise(function(resolve,reject){
		var xhr = new XMLHttpRequest
		xhr.onreadystatechange = function(){
			if (xhr.readyState!==4) return
			if (xhr.status!==200) {
				return reject(new Error("bad status:"+xhr.status))
			}
			resolve(JSON.parse(xhr.responseText))
		}
		xhr.open('GET', path)
		xhr.send()
	})
}
HTMLElement.prototype.append = function(tag,text){
	var e
	if (/^</.test(tag)) {
		e = document.createElement('div')
		e.innerHTML = tag
		e = e.firstChild
	} else {
		e = document.createElement(tag)
	}
	if (text) e.innerHTML = text
	return this.appendChild(e)
}
HTMLElement.prototype.empty = function(tag,text){
	while(this.firstChild) this.removeChild(this.firstChild)
}
var	$  = document.querySelectorAll.bind(document),
	$1 = document.querySelector.bind(document)
NodeList.prototype.__proto__ = Array.prototype
NodeList.prototype.on = function(type,f){
	this.forEach(function(n){
		n.addEventListener(type, f, true)
	})
	return this
}


enrichUserInfo(document.body)
