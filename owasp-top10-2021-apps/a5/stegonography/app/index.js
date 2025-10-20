// Stegonography steganography app
const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const router = express.Router();
const cookieParser = require('cookie-parser');
require("dotenv-safe").load();
const jwt = require('jsonwebtoken');
var mongo = require('mongodb')
const crypto = require('crypto');

// Configures everything needed for the app
app.use(express.static('static'));
app.use('/css', express.static('./css'));
app.use('/js', express.static('./js'));
app.use('/images', express.static('./images'));
app.set('views', 'static/views');
app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true}));
app.use(cookieParser())

// Creates a connection to the database
var port = process.env.MONGO_PORT
var MongoClient = require('mongodb').MongoClient;
var url = "mongodb://db:27017/stego"

// Connect to the database
MongoClient.connect(url, function(err, db) {
    if (err) throw err;
    console.log("Database created!");
    db.close();
});

// Create "users" collection
var url = "mongodb://db:27017/stego"
MongoClient.connect(url, function(err, db) {
    if (err) throw err;
    var dbo = db.db("stego");
    dbo.createCollection("users", function(err, ress) {
        if (err) throw err;
        console.log("Users collection created!");
        db.close();
    })
});

// Add "admin" default user to the database
MongoClient.connect(url, function(err, db) {
    if (err) throw err;
    var dbo = db.db("stego");
    // Create a secure hash for the default admin password instead of storing plaintext
    var salt = crypto.randomBytes(16).toString('hex');
    var ITER = 100000;
    var KEYLEN = 64;
    var DIGEST = 'sha512';
    var hash = crypto.pbkdf2Sync('admin', salt, ITER, KEYLEN, DIGEST).toString('hex');
    var myobj = { username: "admin", password: salt + '$' + hash, iterations: ITER, keylen: KEYLEN, digest: DIGEST };
    dbo.collection("users").insertOne(myobj, function(err, res) {
        if (err) throw err;
        console.log("Admin user added to the database");
        db.close();
    });
});

// User login route, get webpage
router.get("/login", function(req,res) {
    res.render("login.html");
})

// User login route, submit POST request to server
router.post("/login", function(req,res)  {
    var username = req.body.user.name;
    var password = req.body.user.password;
    
    // Verifies user credentials
    function VerifiesUser(callback) {
        MongoClient.connect(url, function(err, db) {
            if (err) throw err;
            var dbo = db.db("stego");
            // Query by username only, then verify password using a secure hash (PBKDF2)
            var query = { username: username };
            dbo.collection("users").find(query).toArray(function(err, result) {
                if (err) throw err;
                db.close();
                if (result.length == 0) {
                    callback('not_found');
                } else {
                    var user = result[0];
                    if (!user.password) {
                        // Missing password field or unexpected format
                        callback('not_found');
                        return;
                    }
                    // Expected stored format: salt$hash (hex)
                    var parts = user.password.split('$');
                    if (parts.length !== 2) {
                        callback('not_found');
                        return;
                    }
                    var salt = parts[0];
                    var storedHash = parts[1];
                    var ITER = user.iterations || 100000;
                    var KEYLEN = user.keylen || 64;
                    var DIGEST = user.digest || 'sha512';

                    var derived = crypto.pbkdf2Sync(password, salt, ITER, KEYLEN, DIGEST).toString('hex');
                    try {
                        var derivedBuf = Buffer.from(derived, 'hex');
                        var storedBuf = Buffer.from(storedHash, 'hex');
                        if (derivedBuf.length !== storedBuf.length) {
                            callback('not_found');
                            return;
                        }
                        if (crypto.timingSafeEqual(derivedBuf, storedBuf)) {
                            callback(user.username);
                        } else {
                            callback('not_found');
                        }
                    } catch(e) {
                        callback('not_found');
                    }
                }
            });
        });
    };

    VerifiesUser((username) => { 
        if (username == "admin") {
            var token = jwt.sign({ username }, process.env.SECRET, {
                expiresIn: 300 // Token expires in 5 minutes
            });
            res.cookie('nodejsSessionToken', token).redirect(301, "/admin");
        } else {
            res.status(500).send('Invalid username or password!').redirect(301, "/logout");
        }
    });
})

// Logout route to deauthorize user session tokens
router.get("/logout", function(req, res) {
    res.status(200).clearCookie('nodejsSessionToken').redirect(301, "/");
});

// Admin maintenance page
router.get("/admin", verifyJWT, (req, res, next) => {
    res.status(200).render("admin.html");
});

// Change password route
router.get("/changepassword", verifyJWT, function(req, res, next) {
    // Code to change user password in the database
})

// Healthcheck route
router.get("/healthcheck", function(req,res) {
    res.send("WORKING");
})

// Main page
router.get("/", function(req,res) {
    res.render("index.html")
})

// Returns the error web-page if none other is found
app.use('/', router);
app.use(function(req, res, next) {
    res.status(404).render("error.html")
});
// Listen on port 10006
app.listen(10006, () => {
    console.log("Server running on port 10006!");
})

// Verifies the JWT token
function verifyJWT(req, res, next){
    var token = req.cookies.nodejsSessionToken;
    if (!token) return res.status(401).send({auth: false, message: 'No token provided'});

    jwt.verify(token, process.env.SECRET, function(err, decoded) {
        if (err) return res.status(500).send({ auth: false, message: 'Failed to authenticate token.' });
        
        req.userId = decoded.id;
        next();
      });
}