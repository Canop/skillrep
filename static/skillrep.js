function fetchUsersPage(page){
	$.getJSON('users')
	.then(function fillTbl(resp){
		console.log('received', resp);
		$('#users tbody').empty().append(resp.Users.map(function(u){
			return $('<tr>')
			.append($('<td>').text(u.Rank))
			.append($('<td>').text(u.Id))
			.append($('<td>').text(u.Name))
			.append($('<td>').text(u.SkillRep))
		}));	
	});
}

$(function(){
	fetchUsersPage(1);
});
