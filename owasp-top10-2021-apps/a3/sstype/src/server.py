import tornado.template
import tornado.ioloop
import tornado.web
import os

with open(os.path.join(os.path.dirname(__file__)) + "/public/index.html", 'r') as f:
    tmpl = f.read()

# Compile the template once at startup; user input is passed as a render variable,
# never injected into the template source.
t = tornado.template.Template(tmpl)

class MainHandler(tornado.web.RequestHandler):

    def get(self):
        name = self.get_argument('name', '')
        self.write(t.generate(name=name))

application = tornado.web.Application([
    (r"/", MainHandler),
    (r"/images/(.*)",tornado.web.StaticFileHandler, {"path": os.path.join(os.path.dirname(__file__)) + "/images/"},),
], debug=False)

if __name__ == '__main__':
    application.listen(10001)
    tornado.ioloop.IOLoop.instance().start()