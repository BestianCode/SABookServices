function checkPassword(form) {
	pw1 = form.password_x1.value;
	pw2 = form.password_x2.value;

	if (pw1 != pw2) {
		alert("Введенные пароли не совпадают!")
		return false;
	}
	if (pw1.length<5) {
		alert("Пароль не должен быть короче 5 символов!")
		return false;
	}

	return true;
}
function checkPasswordAndLogin(form) {
	pw1 = form.password_x1.value;
	pw2 = form.password_x2.value;
	login = form.login.value;

	if (pw1 != pw2) {
		alert("Введенные пароли не совпадают!")
		return false;
	}
	if (pw1.length<5) {
		alert("Пароль не должен быть короче 5 символов!")
		return false;
	}
	if (login.length<5) {
		alert("Логин не должен быть короче 5 символов!")
		return false;
	}

	login_cut = login.replace(/[A-Za-z0-9\.\_\-\@]/g, "")
	if (login_cut.length>0) {
		alert("Логин может содержать большие и маленькие латинские буквы, арабские цифры и символы: . _ @")

		return false;
	}

	return true;
}
function checkDelete(form) {
	check_del = form.check_del.value;

	if (check_del != "удалить нах") {
		alert("Введите правильно то, что просят...")
		return false;
	}
	return true;
}
