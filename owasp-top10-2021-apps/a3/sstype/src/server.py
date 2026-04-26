import tornado.template
import tornado.ioloop
import tornado.web
import os

TEMPLATE = open(os.path.join(os.path.dirname(__file__)) + "/public/index.html", 'r').readlines()

tmpl = ''
for t in TEMPLATE:
    	tmpl += t

class MainHandler(tornado.web.RequestHandler):

    def get(self):
        name = self.get_argument('name', '')
        safe_name = tornado.escape.xhtml_escape(name)
        template_data = tmpl.replace("NAMEHERE", safe_name)
        t = tornado.template.Template(template_data)
        self.write(t.generate(name=safe_name))

application = tornado.web.Application([
    (r"/", MainHandler),
    (r"/images/(.*)",tornado.web.StaticFileHandler, {"path": os.path.join(os.path.dirname(__file__)) + "/images/"},),
], debug=False)

if __name__ == '__main__':
    application.listen(10001)
    tornado.ioloop.IOLoop.instance().start()