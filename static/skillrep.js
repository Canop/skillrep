
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
			n.addEventListener(type, f)
		})
		return this
	}

	var page = 0
	function fetchUsersPage(){
		fetch('users',{page:page, search:$1('#search').value})
		.then(function fillTbl(resp){
			console.log('received', resp)
			if (resp.Error) {
				console.log("error:", resp.Error)
				return
			}
			var tbody = $1("#users tbody")
			tbody.empty()
			resp.Users.forEach(function(u){
				var tr = tbody.append('tr')
				tr.append('td', u.Rank)
				tr.append('td', '<img width=40px height=40px src="'+u.Profile+'">')
				tr.append('td', u.Name)
				tr.append('td', u.Upvotes)
				tr.append('td', u.Accepts)
				tr.append('td', u.Score)
			})
			$1('#nb-questions').textContent = resp.DBStats.NbQuestions
			$1('#nb-answers').textContent = resp.DBStats.NbAnswers
			$1('#max-question-date').textContent = new Date(
				resp.DBStats.MaxQuestionCreationDate*1000
			)
		})
	}

	fetchUsersPage()
	$('.tbl-first').on('click', function(){
		page = 0
		fetchUsersPage()
	})
	$('.tbl-prev').on('click', function(){
		page--
		fetchUsersPage()
	})
	$('.tbl-next').on('click', function(){
		page++
		fetchUsersPage()
	})
	$('#search').on('keyup', function(){
		this.className = this.className.replace(/\binvalid\b/,'')
		try {
			new RegExp(this.value)
			fetchUsersPage()
		} catch(e) {
			this.className += ' invalid'
		}
	});
	var aboutDisplayed = false;
	$('#about-opener').on('click', function(){
		var about = $1('#about')
		about.style.display = about.style.display!=='block'?'block':'none'
	});

}()
