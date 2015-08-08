
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

	var page = 0
	function commated(n){
		return (''+n).replace(/(^\d{1,3}|\d{3})(?=(?:\d{3})+(?:$|\.))/g, '$1,')
	}
	function fetchDBStats(){
		fetch('dbstats',{})
		.then(function fillTbl(resp){
			console.log('received', resp)
			if (resp.Error) {
				console.log("error:", resp.Error)
				return
			}
			$1('#nb-players').textContent = "about " + commated(resp.DBStats.NbPlayers)
			$1('#nb-questions').textContent = "about " + commated(resp.DBStats.NbQuestions)
			$1('#nb-answers').textContent = "about " + commated(resp.DBStats.NbAnswers)
			$1('#max-question-date').textContent = new Date(
				resp.DBStats.MaxQuestionCreationDate*1000
			).toString().replace(/:[^:]*$/,'')
		})
	}
	var currentQuery, postponedQuery
	function fetchQuery(q){
		currentQuery = q
		$1('#wait').className='on'
		fetch('users', q)
		.then(function fillTbl(resp){
			$1('#wait').className='off'
			console.log('received', resp)
			currentQuery = null
			if (resp.Error) {
				console.log("error:", resp.Error)
				return
			}
			if (resp.Users) {
				var tbody = $1("#users tbody")
				tbody.empty()
				resp.Users.forEach(function(u){
					var tr = tbody.append('tr')
					tr.append('td', u.Rank)
					tr.append('td', '<img width=40px height=40px src="'+u.Profile+'">')
					tr.append('td', u.Name)
					tr.append('td', commated(u.SkillRep))
					tr.onclick = function(){
						showUser(u)
					}
				})
			}
			if (postponedQuery) {
				fetchQuery(postponedQuery)
				postponedQuery = null
			}
		})
	}
	function fetchUsersPage(){
		query = {page:page, search:$1('#search').value}
		if (currentQuery) {
			console.log('one at a time, sorry') // search is too slow...
			postponedQuery = query
			return
		}
		fetchQuery(query)
	}
	function hideUser(){
		var $u = $1('#user-wrapper')
		if ($u) $u.remove()
	}
	function showUser(u){
		hideUser()
		var $uw = $1("body").append('div')
		$uw.id = "user-wrapper"
		var $u = $uw.append('div')
		$u.id = "user"
		$u.append('h1', u.Name)
		$u.append('<img src="'+u.Profile+'">')
		$u.append('div', '<span>Rank:</span>').append('<span id=user-rank>', u.Rank)
		$u.append('div', '<span>Skill Rep:</span>').append('<span id=user-skill-rep>', commated(u.SkillRep))
		$u.append('div', '<span>Accepted Answers:</span>').append('<span id=user-accepts>counting...')
		$u.append('div', '<a href=https://stackoverflow.com/users/'+u.Id+'>Stack Overflow Profile</a>')
		fetch("users/"+u.Id+"?accepts=true")
		.then(function(ur){
			$1('#user-accepts').innerHTML = commated(ur.User.Accepts)
		})
	}
	$('body').on('click', hideUser).on('keyup', hideUser)

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
	$1('#username-clear').style.visibility = 'hidden'
	$('#username-clear').on('click', function(){
		$1('#search').value = ''
		page = 0
		this.style.visibility = 'hidden'
		fetchUsersPage()
	})
	var keyTimer
	$('#search').on('keyup', function(){
		this.className = this.className.replace(/\binvalid\b/,'')
		page = 0
		$1('#username-clear').style.visibility = this.value ? 'visible' : 'hidden'
		try {
			new RegExp(this.value)
			clearTimeout(keyTimer)
			keyTimer = setTimeout(fetchUsersPage, 500)
		} catch(e) {
			this.className += ' invalid'
		}
	})
	var uidmatch = location.search.match(/(?:\?|&)user=(\d+)(?:$|&)/)
	if (uidmatch) {
		fetch("users/"+uidmatch[1])
		.then(function(ur){
			if (ur.Error || !ur.User.Name) {
				console.log("bad user", ur)
				return
			}
			showUser(ur.User)
		})
	}
	$('.tab').on('click', function(){
		var selected = /selected/.test(this.className)
		$('.tab').forEach(function(tab){
			tab.className =  'tab'
		})
		$('.tab-page').forEach(function(page){
			page.className = 'tab-page'
		})
		if (selected) return
		this.className = 'tab selected'
		if (this.dataset.page==='about') fetchDBStats()
		document.getElementById(this.dataset.page).className = 'tab-page selected'
	})
}()
