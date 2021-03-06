package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Para utilizar a API utilize /pessoas ou /pessoas/{id} :D")
}

func GetAllPessoas(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:luizgostaderuby@tcp(138.197.3.102:3306)/dojo")
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT p.id, p.nome, p.cpf, p.senha, e.email FROM pessoas p JOIN emails e on p.id = e.pessoa_id")
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}

	pessoas := Pessoas{}

	for rows.Next() {
		var id int
		var nome string
		var cpf string
		var senha string
		var email string

		err = rows.Scan(&id, &nome, &cpf, &senha, &email)

		pessoa := Pessoa{ID: id, Nome: nome, Email: email, CPF: cpf, Senha: senha}

		pessoas = append(pessoas, pessoa)

	}

	if err := json.NewEncoder(w).Encode(pessoas); err != nil {
		panic(err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func GetPessoaID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	idPessoa, err := strconv.Atoi(vars["pessoaId"])

	if err != nil {
		fmt.Fprintln(w, "ID inválido")
		panic(err.Error())
	}

	db, err := sql.Open("mysql", "root:luizgostaderuby@tcp(138.197.3.102:3306)/dojo")
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT p.id, p.nome, p.cpf, p.senha, e.email FROM pessoas p JOIN emails e on p.id = e.pessoa_id WHERE p.id = " + strconv.Itoa(idPessoa))
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}

	if rows.Next() {
		var id int
		var nome string
		var cpf string
		var senha string
		var email string

		err = rows.Scan(&id, &nome, &cpf, &senha, &email)

		pessoa := Pessoa{ID: id, Nome: nome, Email: email, CPF: cpf, Senha: senha}

		if err := json.NewEncoder(w).Encode(pessoa); err != nil {
			panic(err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Não existe nemhum usuário com o ID "+strconv.Itoa(idPessoa))

	}

}

func DeletePessoaID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	idPessoa, err := strconv.Atoi(vars["pessoaId"])

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintln(w, "ID inválido")
		panic(err.Error())
	}

	db, err := sql.Open("mysql", "root:luizgostaderuby@tcp(138.197.3.102:3306)/dojo")
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}
	defer db.Close()

	delPessoa, err := db.Prepare("DELETE FROM pessoas WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erro interno tente novamente")
	} else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "Excluido com sucesso")
	}

	delPessoa.Exec(idPessoa)
	defer db.Close()

}

func CriarPessoa(w http.ResponseWriter, r *http.Request) {

	// vars := mux.Vars(r)
	w.WriteHeader(http.StatusPartialContent)
	fmt.Fprintf(w, r.FormValue("nome"))

	db, err := sql.Open("mysql", "root:luizgostaderuby@tcp(138.197.3.102:3306)/dojo")
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}
	defer db.Close()

	criarPessoa, err := db.Prepare("INSERT INTO pessoas (nome, cpf, senha) VALUES(?, ?, ?);")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erro interno tente novamente")
	}

	criarEmail, err := db.Prepare("INSERT INTO emails (email, pessoa_id) VALUES(?, LAST_INSERT_ID());")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erro interno tente novamente")
	}

	criarPessoa.Exec(r.FormValue("nome"), r.FormValue("cpf"), r.FormValue("senha"))
	criarEmail.Exec(r.FormValue("email"))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, " Usuário Inserido")

}

func AlterarSenhaPessoa(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idPessoa, err := strconv.Atoi(vars["pessoaId"])

	db, err := sql.Open("mysql", "root:luizgostaderuby@tcp(138.197.3.102:3306)/dojo")
	if err != nil {
		fmt.Print("DEU RUIM")
		panic(err.Error())
	}
	defer db.Close()

	alterarSenha, err := db.Prepare("UPDATE pessoas SET senha=? WHERE id=?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Erro interno tente novamente")
	}

	alterarSenha.Exec(r.FormValue("senha"), idPessoa)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Senha Alterada")

}
