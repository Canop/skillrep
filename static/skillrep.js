
!function(){
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
		var e = document.createElement(tag)
		if (text) e.textContent = text
		return this.appendChild(e)
	}
	var $ = document.querySelectorAll.bind(document)
	NodeList.prototype.__proto__ = Array.prototype
	NodeList.prototype.on = function(type,f){
		this.forEach(function(n){
			n.addEventListener(type, f)
		})
		return this
	}

	function fetchUsersPage(page){
		fetch('users')
		.then(function fillTbl(resp){
			console.log('received', resp)
			var tbody = document.querySelector("#users tbody")
			resp.Users.forEach(function(u){
				var tr = tbody.append('tr')
				tr.append('td', u.Rank)
				tr.append('td', u.Id)
				tr.append('td', u.Name)
				tr.append('td', u.SkillRep)
			})
		})
	}

	fetchUsersPage(1)
}()
