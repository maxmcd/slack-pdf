const app = require("express")();
const multer  = require('multer')
var upload = multer()

const spawn = require('child_process').spawn;

app.post('/', upload.single('file'), function (req, res, next) {
  res.set("Content-Disposition", `attachment; filename=${req.file.originalname}.pdf`)
  let child = spawn('cat')
  child.stdin.write(req.file.buffer)
  child.stdin.end()
  child.stdout.on('data', (data) => {
    res.write(data)
  });
  child.on('close', (code) => {
    res.end()
  });
})

app.listen(8080, function () {
  console.log('up 8080')
})