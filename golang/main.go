package main

import (
	
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	/* "./encriptacion/cookie"
	"github.com/pacomisval/Gnommo_api/GO/encriptacion/token"
	"github.com/pacomisval/Gnommo_api/GO/encriptacion/uploadFile"
	"github.com/pacomisval/Gnommo_api/GO/encriptacion/recoveryPass"
	"github.com/pacomisval/Gnommo_api/GO/libreria/autor"
	"github.com/pacomisval/Gnommo_api/GO/libreria/libro"
	"github.com/pacomisval/Gnommo_api/GO/libreria/usuario" */
)

type Libro struct {
	Id          string `json:"id"`
	Nombre      string `json:"nombre"`
	Isbn        string `json:"isbn"`
	Genero      string `json:"genero"`
	Descripcion string `json:"descripcion"`	
	Portada     string `json:"portada"`
	IdAutor     string `json:"idAutor"`

	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type Autor struct {
	Id        		  string `json:"id"`
	FirstName 		  string `json:"first_name"`
	LastName  		  string `json:"last_name"`
	Nacionalidad 	  string `json:"nacionalidad"`
	FechaNacimiento   string `json:"fechaNacimiento"`
}

type Usuario struct {
	Id       string `json:"id"`
	Nombre   string `json:"nombre"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Rol      string `json:"rol"`
	Tok      string `json:"tok"`
	Codigo   string `json:"codigo"`
}
type Claims struct {
	//Id     uint
	Nombre string
	Email  string
	*jwt.StandardClaims
}
type Value struct {
	Id     string `json:"id"`
	Nombre string `json:"Nombre"`
	Rol    string `json:"rol"`
	Token  string `json:"token"`
}

var db *sql.DB
var err error

const maxUploadSize = 200 * 1024 // 100 KB
const uploadPath = "./src/assets/images/book"

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/newlibrary")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/libros", getLibros).Methods("GET")
	router.HandleFunc("/api/libros/autor/{id}", getLibrosByAutor).Methods("GET")
	router.HandleFunc("/api/libros/all", getAll).Methods("GET")
	router.HandleFunc("/api/libros/{id}", getLibro).Methods("GET")
	router.HandleFunc("/api/libros", postLibro).Methods("POST")
	router.HandleFunc("/api/libros/{id}", putLibro).Methods("PUT")
	router.HandleFunc("/api/libros/{id}", deleteLibro).Methods("DELETE")

	router.HandleFunc("/api/autores", getAutores).Methods("GET")
	router.HandleFunc("/api/autores/{id}", getAutor).Methods("GET")
	router.HandleFunc("/api/autores", postAutor).Methods("POST")
	router.HandleFunc("/api/autores/{id}", putAutor).Methods("PUT")
	router.HandleFunc("/api/autores/{id}", deleteAutor).Methods("DELETE")

	router.HandleFunc("/api/usuarios", getUsuarios).Methods("GET")
	router.HandleFunc("/api/usuarios/{id}", getUsuario).Methods("GET")
	router.HandleFunc("/api/usuarios", postUsuario).Methods("POST")
	router.HandleFunc("/api/usuarios/{id}", putUsuario).Methods("PUT")
	router.HandleFunc("/api/usuarios/{id}", deleteUsuario).Methods("DELETE")

	router.HandleFunc("/api/registro", postUsuario).Methods("POST")
	router.HandleFunc("/api/login", login).Methods("POST")

	router.HandleFunc("/api/recoveryPass1", recuperarPass).Methods("POST")
	router.HandleFunc("/api/recoveryPass2", verificarCodigo).Methods("POST")
	router.HandleFunc("/api/recoveryPass3", nuevoPassword).Methods("POST")

	router.HandleFunc("/api/upload", uploadFileHandler).Methods("POST")

	log.Print("Server started on localhost:8000 ..........")

	log.Print("Otra linea mas ......................................")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Accept", "Accept-Language"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:4200"}))(router)))

}

////////////////////////////////////// INICIO ENCRIPTACION /////////////////////////////////////
 
func encriptarPass(pass string, clave string) string {
	hashMD5 := MD5Hash(pass)
	hashHMAC := HMACHash(hashMD5, clave)
	return hashHMAC
}

func MD5Hash(pass string) string {
	hash := md5.New()
	hash.Write([]byte(pass))
	return hex.EncodeToString(hash.Sum(nil))
}

func HMACHash(pass string, clave string) string {
	hash := hmac.New(sha256.New, []byte(clave))
	io.WriteString(hash, pass)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func validarPass(passInput, passwordDB string) bool {

	bytePassDB := []byte(passwordDB)

	byteInput := []byte(passInput)

	resHMAC := hmac.Equal(byteInput, bytePassDB)

	return resHMAC
}

////////////////////////////////////////////////////////////////////////////////////////
///////////////////  PASO 1  RECUPERAR CONTRASEÑA  //////////////////////////////////////////

func recuperarPass(w http.ResponseWriter, r *http.Request) {
	var value bool
	var id string

	fmt.Println("Dentro de recuperarPass")

	useru := &Claims{}
	err := json.NewDecoder(r.Body).Decode(useru)

	fmt.Println("Valor de err: ", err)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
	}

	value, id = findUsuarioByEmail(useru.Email)

	fmt.Println("Valor de user.Nombre: ", useru.Nombre)
	fmt.Println("Valor de user.Email: ", useru.Email)
	fmt.Println("Valor de useru: ", useru)
	fmt.Println("")
	fmt.Println("Valor de value en recuperarPass: ", value)
	fmt.Println("Valor de id en recuperarPass: ", id)
	fmt.Println("Valor de w en recuperarPass: ", w)

	if value{

		clave := crearClaveAleatoria(id, useru.Email)
		fmt.Println("valor de clave: ", clave)

		sendMailRecoveryPassword(clave)
		fmt.Println("Se ha enviado un correo")
	}

	json.NewEncoder(w).Encode(value)

}

func findUsuarioByEmail( mail string) ( bool, string) {
	ok := false

	fmt.Println("Valor de mail en findUsuarioByEmail: ", mail)

	var id string
	var email string

	result, err := db.Query("SELECT id, email FROM usuario WHERE email like ?", &mail)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var usuario Usuario
		err := result.Scan(&usuario.Id, &usuario.Email)
		if (err) != nil {
			panic(err.Error())
		}

		id = usuario.Id
		email = usuario.Email
	}

	resultEmail := mail == email

	fmt.Println("")
	fmt.Println("Valor de resultEmail: ", resultEmail)
	fmt.Println("")
	fmt.Println("Valor de email: ", email)
	fmt.Println("Valor de mail: ", mail)
	fmt.Println("Valor de id: ", id)
	fmt.Println("")

	if resultEmail {
		ok = true
		fmt.Println("Valor de ok: ", ok)
	}

	return ok, id
}

func crearClaveAleatoria(id, email string) string {

	expireAt := time.Now().Add(time.Hour * 50).Unix()
	claims := Claims{
		Email:  email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
	}

	Secret := []byte(email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		panic(err.Error())
	}

	c := []rune(tokenString)
	codigo := string(c[124:134])

	fmt.Println("Valor de tokenString: ", tokenString)
	fmt.Println("Valor de codigo: ", codigo)

	guardarToken(codigo, id)

	return codigo
}

func sendMailRecoveryPassword(body string) { // pasar variables con el from, pass, to
	// Entra en esta url: https://support.google.com/accounts/answer/6010255
	// "Hay que habilitar en mi cuenta el modo inseguro"
	// Entra aqui: Si tu cuenta tiene activado el acceso de aplicaciones poco seguras
	// Haz click en el link:  Acceso de aplicaciones poco seguras
	// La opción debe estar así: Permitir el acceso de aplicaciones poco seguras: SÍ

	from := "Aqui tu email"                      // tu email
	pass := "Aqui tu contraseña de gmail"        // tu contraseña de gmail o google
	to := "Aqui el email a quien quieres enviar" // email del usuario que renueva la contraseña

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello Mari. Introduce este código\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, dentro de send")
}

///////////////  FIN PASO 1  RECUPERAR CONTRASEÑA  /////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
///////////////  PASO 2  RECUPERAR CONTRASEÑA  ////////////////////////////////

func verificarCodigo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dentro de verificarCodigo")

	code := &Usuario{}
	err := json.NewDecoder(r.Body).Decode(code)

	fmt.Println("Valor de err: ", err)
	fmt.Println("Valor de code.Codigo: ", code.Codigo)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
	}

	// Aqui hay un error, No se obtiene el registro  ¿Porque ?
	result, err := db.Query("SELECT * FROM usuario WHERE codigo like ?", code.Codigo)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var usuario Usuario
	for result.Next() {

		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if (err) != nil {
			panic(err.Error())
		}
	}
	resultCodigo := usuario.Codigo == code.Codigo
	fmt.Println("Valor de usuario.codigo: ", usuario.Codigo)
	fmt.Println("Valor de resultCodigo: ", resultCodigo)

	id := usuario.Id
	nombre := usuario.Nombre
	passwd := usuario.Password
	email := usuario.Email
	rol := usuario.Rol
	tok := usuario.Tok
	codigo := ""

	stmt, err := db.Prepare("UPDATE usuario SET id = ?, nombre = ?, password = ?, email = ?, rol = ?, tok = ?, codigo = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(id, nombre, passwd, email, rol, tok, codigo, id)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ESTO ES ELIMINAR CODIGO")

	json.NewEncoder(w).Encode(resultCodigo)
}

/////////////////  FIN PASO 2  RECUPERAR CONTRASEÑA  ////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////  PASO 3  RECUPERAR CONTRASEÑA  ///////////////////////////////

func nuevoPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("Dentro de nuevoPassword:")
	var ok = true

	user := &Usuario{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		//var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		ok = false
		json.NewEncoder(w).Encode(ok)
	}
	fmt.Println("nuevoPassword valor de user.Email: ", user.Email)
	fmt.Println("nuevoPassword valor de user.Password: ", user.Password)

	result, err := db.Query("SELECT * FROM usuario WHERE email like ?", user.Email)
	if err != nil {
		//panic(err.Error())
		ok = false
		json.NewEncoder(w).Encode(ok)
	}
	defer result.Close()

	var usuario Usuario
	for result.Next() {

		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if (err) != nil {

		}
	}
	id := usuario.Id
	nombre := usuario.Nombre
	password := usuario.Password
	email := usuario.Email
	rol := usuario.Rol
	tok := usuario.Tok
	codigo := usuario.Codigo

	newPassword := encriptarPass(user.Password, user.Email)
	fmt.Println("nuevoPassword valor de newPassword: ", newPassword)
	fmt.Println("nuevoPassword valor de password: ", password)

	stmt, err := db.Prepare("UPDATE usuario SET id = ?, nombre = ?, password = ?, email = ?, rol = ?, tok = ?, codigo = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())

	}

	_, err = stmt.Exec(id, nombre, newPassword, email, rol, tok, codigo, id)
	if err != nil {
		//panic(err.Error())
		ok = false
		json.NewEncoder(w).Encode(ok)
	}

	//var resp = map[string]interface{}{"status": "0k", "message": "New password created"}
	json.NewEncoder(w).Encode(ok)

	fmt.Println("ESTO ES NUEVA CONTRASEÑA RECUPERADA")

}

////////////////////////////  FIN PASO 3  RECUPERAR CONTRASEÑA  //////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////  LOGIN ///////////////////////////////////////////////////////

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entra en login ")

	user := &Usuario{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	fmt.Println("Va a encontrar usuario: ")

	value := encontrarUsuario(user.Password, user.Email)
	fmt.Println("Vuelve encontrar usuario: ")

	crearCookie(w, r, value)

	fmt.Println("valor de value: ", value)

	json.NewEncoder(w).Encode(value)
	fmt.Println("El valor de W: ", w)

	return
}

//////////////////////////////////////////  TOKEN  ////////////////////////////////////////

func encontrarUsuario(password, email string) Value {
	var user []Usuario
	var ok bool
	var id string
	var nombre string
	var passDB string
	var mail string
	var rol string
	var resultVacio bool
	resultVacio = true

	expireAt := time.Now().Add(time.Minute * 5).Unix()
	result, err := db.Query("SELECT * FROM usuario WHERE email like ?", &email)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		resultVacio = false
		fmt.Println("entra en for result.next")
		var usuario Usuario
		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if err != nil {
			panic(err.Error())
		}
		user = append(user, usuario)
		id = usuario.Id
		nombre = usuario.Nombre
		passDB = usuario.Password
		mail = usuario.Email
		rol = usuario.Rol
	}

	if resultVacio {

		Id := "0"
		Nombre := "error"
		Rol := "NO se ha encontrado el Email"
		Token := "Email not Found"
		value := Value{Id, Nombre, Rol, Token}
		return value
	}
	contra := encriptarPass(password, email)

	ok = validarPass(contra, passDB)
	if !ok {
		Id := "0"
		Nombre := "error"
		Rol := "Contraseña Incorrecta"
		Token := "Error NO coinciden las contraseñas"
		value := Value{Id, Nombre, Rol, Token}
		return value
	}

	claims := Claims{
		Nombre: nombre,
		Email:  mail,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
	}

	Secret := []byte(mail)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		Id := "0"
		Nombre := "error"
		Rol := "NO se ha creado el token"
		Token := "Error en el token"
		value := Value{Id, Nombre, Rol, Token}
		return value
	}

	guardarToken(tokenString, id)

	value := Value{id, nombre, rol, tokenString}
	return value
}

func guardarToken(token, id string) {

	var iddb string
	var nombre string
	var passwd string
	var email string
	var rol string
	var tok string
	var codigo string

	result, err := db.Query("SELECT * FROM usuario WHERE id = ?", &id)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	fmt.Println("dentro de guardar token: ")

	var usuario Usuario

	for result.Next() {
		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if err != nil {
			panic(err.Error())
		}
	}

	length := len(token)
	fmt.Println("El tamaño de length: ", length)
	if length >= 11 {
		iddb = usuario.Id
		nombre = usuario.Nombre
		passwd = usuario.Password
		email = usuario.Email
		rol = usuario.Rol
		tok = token
		codigo = usuario.Codigo
	} else {
		iddb = usuario.Id
		nombre = usuario.Nombre
		passwd = usuario.Password
		email = usuario.Email
		rol = usuario.Rol
		tok = usuario.Tok
		codigo = token
	}
	fmt.Println("Codigo dentro de guardar token: ", token)
	fmt.Println("Usuario.codigo dentro de guardar token: ", usuario.Codigo)

	// Error, No guarda el codigo en la DB, No actualiza el registro,  Nose porque?
	stmt, err := db.Prepare("UPDATE usuario SET id = ?, nombre = ?, password = ?, email = ?, rol = ?, tok = ?, codigo = ?  WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(iddb, nombre, passwd, email, rol, tok, codigo, id)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("UPDATE en guardar token: ", codigo)

	fmt.Println("ESTO ES GUARDAR TOKEN")

}

func recuperarToken(id string) string {
	result, err := db.Query("SELECT tok FROM usuario WHERE id = ?", &id)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var usuario Usuario

	for result.Next() {
		err := result.Scan(&usuario.Tok)
		if err != nil {
			panic(err.Error())
		}
	}

	return usuario.Tok
}

func crearToken(id, nombre, email string) string {

	expireAt := time.Now().Add(time.Minute * 125).Unix()
	claims := Claims{
		Nombre: nombre,
		Email:  email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
	}

	Secret := []byte(email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		panic(err.Error())
	}

	guardarToken(tokenString, id)

	return tokenString
}

func verificarToken(tknStr, SecretKey string) int {
	resp := 0 // ok
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	fmt.Println("Valor del token: ", tknStr)
	fmt.Println("Valor de claims: ", claims)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			resp = 1 // StatusUnauthorized
			return resp
		}
		resp = 2 // StatusBadRequest
	}
	if !tkn.Valid {
		resp = 1
		return resp
	}

	fmt.Println("Dentro de verificarToken: ", resp)
	return resp
}

/////////////////////////////////////////////  COOKIES  //////////////////////////////////////////////////////

func crearCookie(w http.ResponseWriter, r *http.Request, value Value) {

	I := value.Id
	N := value.Nombre
	R := value.Rol
	T := value.Token

	expiration := time.Now().Add(time.Minute * 5)

	cookie1 := &http.Cookie{
		Name:     "tokensiI",
		Value:    I,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: false,
		//SameSite: Lax,
		Secure: false,
	}
	http.SetCookie(w, cookie1)
	r.AddCookie(cookie1)

	cookie2 := &http.Cookie{
		Name:     "tokensiN",
		Value:    N,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: false,
		//SameSite: Lax,
		Secure: false,
	}
	http.SetCookie(w, cookie2)
	r.AddCookie(cookie2)

	cookie3 := &http.Cookie{
		Name:     "tokensiR",
		Value:    R,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: false,
		//SameSite: Lax,
		Secure: false,
	}
	http.SetCookie(w, cookie3)
	r.AddCookie(cookie3)

	cookie4 := &http.Cookie{
		Name:     "tokensiT",
		Value:    T,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: false,
		//SameSite: Lax,
		Secure: false,
	}
	http.SetCookie(w, cookie4)
	r.AddCookie(cookie4)
}

func eliminarCookie(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(time.Minute - 1)

	cookie1 := &http.Cookie{
		Name:    "tokensiI",
		Value:   "",
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, cookie1)
	r.AddCookie(cookie1)

	cookie2 := &http.Cookie{
		Name:    "tokensiN",
		Value:   "",
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, cookie2)
	r.AddCookie(cookie2)

	cookie3 := &http.Cookie{
		Name:    "tokensiR",
		Value:   "",
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, cookie3)
	r.AddCookie(cookie3)

	cookie4 := &http.Cookie{
		Name:    "tokensiT",
		Value:   "",
		Path:    "/",
		Expires: expiration,
	}
	http.SetCookie(w, cookie4)
	r.AddCookie(cookie4)
}

func verificarCookies(w http.ResponseWriter, r *http.Request) int {
	c, err := r.Cookie("tokensiT")
	if err != nil {
		fmt.Println("ERROR cookie tokensiT: ", err)
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Error de StatusUnauthorized")
			return 1
		}
	}
	tknStr := c.Value

	d, err := r.Cookie("tokensiI")
	if err != nil {
		fmt.Println("ERROR cookie tokensiI: ", err)
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Error de StatusUnauthorized")
			return 1

		}
	}
	idStr := d.Value

	result, err := db.Query("SELECT * FROM usuario WHERE id = ?", idStr)
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var usuario Usuario

	for result.Next() {

		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if err != nil {
			panic(err.Error())
		}
	}
	secret := usuario.Email

	fmt.Println("Id en putLibro: ", idStr)
	fmt.Println("Token en putLibro: ", tknStr)
	fmt.Println("Valor de secret: ", secret)

	resp := verificarToken(tknStr, secret)

	if resp == 0 {
		token := crearToken(usuario.Id, usuario.Nombre, usuario.Email)

		var value Value
		value.Id = usuario.Id
		value.Nombre = usuario.Nombre
		value.Rol = usuario.Rol
		value.Token = token

		eliminarCookie(w, r)
		crearCookie(w, r, value)
	}

	fmt.Println("Esto es la respuesta de resp en PutLibro: ", resp)

	return resp
}

///////////////////////////////// FIN ENCRIPTACION ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
//////////////////////////  UPLOAD FILES ///////////////////////////////////////

func uploadFileHandler(w http.ResponseWriter, r *http.Request ) {

	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, nil)
		fmt.Println("valor de t: ", t)
		return
	}
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Printf("Could not parse multipart form: %v\n", err)
		renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
		return
	}

	// parse and validate file and post parameters
	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get and print out file size
	fileSize := fileHeader.Size
	fmt.Printf("File size (bytes): %v\n", fileSize)
	// validate file size
	if fileSize > maxUploadSize {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg":
		break
	case "image/jpg":
		break
	case "image/gif":
		break
	case "image/png":
		break
	case "application/pdf":
		break
	default:
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return
	}
	fileName := fileHeader.Filename           //randToken(12)   // archivo
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return
	}
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	fmt.Printf("FileType: %s, FilePath: %s, FileName: %s\n", detectedFileType, newPath, fileName)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		fmt.Println("valor de w1: ", w)
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		fmt.Println("valor de w2: ", w)
		return
	}
	fmt.Println("Todo ha ido bien. Has llegado al final !!")
	//w.Write([]byte("SUCCESS"))

	json.NewEncoder(w).Encode(fileName)
}                                          

func renderError(w http.ResponseWriter, message string, statusCode int) {

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {

	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
} 

//////////////////////////////  FIN UPLOAD FILES  /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////// INICIO API LIBROS ////////////////////////////

//////////////////// GET LIBROS ////////////////////////
 func getLibros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var libros []Libro

	result, err := db.Query("SELECT * FROM libro")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var libro Libro
		err := result.Scan(&libro.Id, &libro.Nombre, &libro.Isbn, &libro.Genero, &libro.Descripcion, &libro.Portada,&libro.IdAutor )
		if err != nil {
			panic(err.Error())
		}
		libros = append(libros, libro)
	}

	json.NewEncoder(w).Encode(libros)

	fmt.Println("ESTO ES GET LIBROS")
}

/////////////// GET LIBROS POR AUTOR ///////////////////////////
func getLibrosByAutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var libros []Libro

	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM libro WHERE idAutor = ?", params["id"])
	if err != nil {
		panic(err.Error())
	
	}
	defer result.Close()

	for result.Next() {
		var libro Libro
		err := result.Scan(&libro.Id, &libro.Nombre, &libro.Isbn, &libro.Genero, &libro.Descripcion, &libro.Portada,&libro.IdAutor)
		if err != nil {
			panic(err.Error())
		}
		libros = append(libros, libro)
	}
	json.NewEncoder(w).Encode(libros)

	fmt.Println("ESTO ES GET LIBRO POR AUTOR")
}

/////////////////////// GET TODO ////////////////////
func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var libros []Libro

	result, err := db.Query("SELECT b.id, b.nombre, b.isbn, b.genero, b.descripcion, b.portada, b.idAutor, a.first_name, a.last_name FROM libro b INNER JOIN autor a ON b.idAutor = a.id ORDER BY genero DESC")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var libro Libro

		err := result.Scan(&libro.Id, &libro.Nombre, &libro.Isbn, &libro.Genero, &libro.Descripcion, &libro.Portada, &libro.IdAutor, &libro.FirstName, &libro.LastName)
		if err != nil {
			panic(err.Error())
		}
		libros = append(libros, libro)
	}
	json.NewEncoder(w).Encode(libros)

	fmt.Println("ESTO ES GET ALL")
}

/////////////// GET LIBRO POR ID /////////////////////
func getLibro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM libro WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var libro Libro

	for result.Next() {
		err := result.Scan(&libro.Id, &libro.Nombre, &libro.Isbn, &libro.Genero, &libro.Descripcion, &libro.Portada, &libro.IdAutor)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(libro)

	fmt.Println("ESTO ES GET BY ID")
}

//////////////////// POST LIBRO ///////////////////////
func postLibro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	///////////////////////////////////////////////
	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}
	///////////////////////////////////////////////

	/* uploadFileHandler(w, r)
	fmt.Println("EL VALOR DEL FICHERO EN POST LIBRO: " , fichero ) */

	//stmt, err := db.Prepare("INSERT INTO books(id, nombre, isbn, idAutor) VALUES(?,?,?,?)")
	stmt, err := db.Prepare("INSERT INTO libro(id, nombre, isbn, genero, descripcion, portada, idAutor) VALUES(?,?,?,?,?,?,?)")

	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	// archivoaux := body.isbn
	// archivo := archivoaux + "."

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	id := keyVal["id"]
	nombre := keyVal["nombre"]
	isbn := keyVal["isbn"]
	genero := keyVal["genero"]
	descripcion := keyVal["descripcion"]
	portada := keyVal["portada"]
	idAutor := keyVal["id_author"] //mirar si falla FK es idAutor

	//fmt.Println("inserta LIBRO:", id, nombre, isbn, idAutor)

	_, err = stmt.Exec(&id, &nombre, &isbn, &genero, &descripcion, &portada, &idAutor)

	if err != nil {
		panic(err.Error())
	}
	//	fmt.Fprintf(w, "Se a añadido un nuevo libro")

	fmt.Println("ESTO ES POST LIBRO")
}

///////////////////// PUT LIBRO /////////////////////
func putLibro(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	///////////////////////////////////////////////

	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}

	///////////////////////////////////////////////
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE libro SET id = ?, nombre = ?, isbn = ?, genero = ?, descripcion = ?, portada = ?, idAutor = ?  WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	claveValor := make(map[string]string)
	json.Unmarshal(body, &claveValor)

	newId := claveValor["id"]
	newNombre := claveValor["nombre"]
	newIsbn := claveValor["isbn"]
	newGenero := claveValor["genero"]
	newDescripcion := claveValor["descripcion"]
	newPortada := claveValor["portada"]
	newIdAutor := claveValor["idAutor"]
	

	_, err = stmt.Exec(&newId, &newNombre, &newIsbn, &newGenero, &newDescripcion, &newPortada, &newIdAutor, params["id"])
	if err != nil {
		panic(err.Error())
	}

	// fmt.Fprintf(w, "El registro con Id %s se ha actualizado correctamente", params["id"])
	fmt.Println("ESTO ES PUT LIBRO")
}

////////////////// DELETE LIBRO /////////////////////////
func deleteLibro(w http.ResponseWriter, r *http.Request) {

	///////////////////////////////////////////////////////////
	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}
	////////////////////////////////////////////////////////////////

	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM libro WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	//	fmt.Fprintf(w, "Se ha eliminado el libro con Id %s", params["id"])
	fmt.Println("ESTO ES DELETE LIBRO")
}

//////////////////////////////////// FIN API LIBROS ///////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////// INICIO API AUTORES ///////////////////////////

//////////////////////// GET AUTORES ///////////////////
func getAutores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var autores []Autor

	result, err := db.Query("SELECT * FROM autor")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var autor Autor
		err := result.Scan(&autor.Id, &autor.FirstName, &autor.LastName, &autor.Nacionalidad, &autor.FechaNacimiento)
		if err != nil {
			panic(err.Error())
		}
		autores = append(autores, autor)
	}
	json.NewEncoder(w).Encode(autores)

	fmt.Println("ESTOS SON TODOS LOS AUTORES")
}


/////////////////// GET AUTOR POR ID //////////////////
func getAutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM autor WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var autor Autor

	for result.Next() {
		err := result.Scan(&autor.Id, &autor.FirstName, &autor.LastName, &autor.Nacionalidad, &autor.FechaNacimiento)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(autor)

	fmt.Println("ESTO ES GET AUTOR POR ID")

}

/////////////////////// POST AUTOR ////////////////////
func postAutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	///////////////////////////////////////////////

	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}

	///////////////////////////////////////////////

	stmt, err := db.Prepare("INSERT INTO autor(id, first_name, last_name, nacionalidad, fechaNacimiento) VALUES (?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	clave := make(map[string]string)
	json.Unmarshal(body, &clave)

	id := clave["id"]
	firstName := clave["first_name"]
	lastName := clave["last_name"]
	nacionalidad := clave["nacionalidad"]
	fechaNacimiento := clave["fechaNacimiento"]

	_, err = stmt.Exec(&id, &firstName, &lastName, &nacionalidad, &fechaNacimiento)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ESTO ES POST AUTOR")
}

////////////////////// PUT AUTOR //////////////////////
func putAutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	///////////////////////////////////////////////

	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}

	///////////////////////////////////////////////
	fmt.Println("ESTO ES PUT AUTOR 1")

	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE autor SET first_name = ?, last_name = ?, nacionalidad = ?,fechaNacimiento = ?  WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	key := make(map[string]string)
	json.Unmarshal(body, &key)

	//nuevoId := key["idAutor"]
	nuevoFirstName := key["first_name"]
	nuevoLastName := key["last_name"]
	nuevaNacionalidad := key["nacionalidad"]
	nuevaFechaNacimiento := key["fechaNacimiento"]

	_, err = stmt.Exec(&nuevoFirstName, &nuevoLastName, &nuevaNacionalidad, &nuevaFechaNacimiento, params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ESTO ES UPDATE AUTOR")
}

//////////////////// DELETE AUTOR ////////////////////////
func deleteAutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	///////////////////////////////////////////////

	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}

	///////////////////////////////////////////////

	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM autor WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ESTO ES DELETE AUTOR")
}

///////////////////////BUSCAR Si existe Autor POR NOMBRE ////////////////////
func buscaAutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("entra en buscar autor")

	//params := mux.Vars(r)
	fmt.Println(r.Body)

	result, err := db.Prepare("SELECT * FROM autor WHERE first_name =? and last_name=?")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	clave := make(map[string]string)
	json.Unmarshal(body, &clave)

	firstName := clave["first_name"]
	lastName := clave["last_name"]

	_, err = result.Exec(&firstName, &lastName)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(result)

}

//////////////////////////////// FIN API AUTORES //////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// INICIO  API USUARIOS //////////////////////////////

////////////////// GET USUARIOS /////////////////////////
func getUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usuarios []Usuario

	result, err := db.Query("SELECT * FROM usuario")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var usuario Usuario
		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if err != nil {
			panic(err.Error())
		}
		usuarios = append(usuarios, usuario)
	}
	json.NewEncoder(w).Encode(usuarios)

	fmt.Println("ESTO ES GET USUARIOS")
}

/////////////// GET USUARIO POR ID //////////////////////
func getUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM usuario WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var usuario Usuario

	for result.Next() {
		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(usuario)

	fmt.Println("ESTO ES USUARIO PO ID")
}

/////////////////// POST USUARIO ////////////////////////

func postUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO usuario(id, nombre, password, email, rol, tok, codigo) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	key := make(map[string]string)
	json.Unmarshal(body, &key)

	id := key["id"]
	nombre := key["nombre"]
	password := key["password"]
	email := key["email"]
	rol := key["rol"]
	tok := key["tok"]
	codigo := key["tok"]

	fmt.Println("valor de nombre: ", nombre)
	fmt.Println("valor de password: ", password)
	fmt.Println("valor de email: ", email)

	pass := encriptarPass(password, email)
	tok = crearToken(id, nombre, email)

	_, err = stmt.Exec(&id, &nombre, &pass, &email, &rol, &tok, &codigo)
	if err != nil {
		panic(err.Error())
	}

	result, err := db.Query("SELECT * FROM usuario WHERE email like ?", &email)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var usuario Usuario

	for result.Next() {
		err := result.Scan(&usuario.Id, &usuario.Nombre, &usuario.Password, &usuario.Email, &usuario.Rol, &usuario.Tok, &usuario.Codigo)
		if err != nil {
			panic(err.Error())
		}
	}

	value := Value{usuario.Id, usuario.Nombre, usuario.Rol, usuario.Tok}

	json.NewEncoder(w).Encode(value)

	fmt.Println("ESTO ES POST USUARIO")
}

///////////////////// PUT USUARIO ///////////////////////
func putUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE usuario SET id = ?, nombre = ?, password = ?, email = ?, rol = ?, tok = ?, codigo = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	key := make(map[string]string)
	json.Unmarshal(body, &key)

	newId := key["id"]
	newNombre := key["nombre"]
	newPassword := key["password"]
	newEmail := key["email"]
	newRol := key["rol"]
	newTok := key["tok"]
	newCodigo := key["codigo"]

	_, err = stmt.Exec(&newId, &newNombre, &newPassword, &newEmail, &newRol, &newTok, &newCodigo, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Se ha actualizado el usuario %s correctamente", params["id"])
	fmt.Println("ESTO ES UPDATE USUARIOS")
}

///////////////////// DELETE USUARIO //////////////////////
func deleteUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	///////////////////////////////////////////////

	resp := verificarCookies(w, r)

	if resp != 0 {
		return
	}

	///////////////////////////////////////////////

	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM usuario WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Se ha eliminado el usuario %s correctamente ", params["id"])
	fmt.Println("ESTO ES DELETE USUARIO")
}
 
////////////////////////////////// FIN API USUARIOS //////////////////////////////
