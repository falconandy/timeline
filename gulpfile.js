var gulp  = require('gulp'),
    util = require('gulp-util'),
    concat = require('gulp-concat'),
    zip = require('gulp-zip'),
    runSequence = require('run-sequence'),
    del = require('del'),
    argv = require('yargs').argv,
    exec = require('child_process').exec,
    os = require('os');

var pkg = require('./package.json');
var platform = os.platform()

gulp.task('default', ['build', 'watch']);

gulp.task('build', function(callback) {
  runSequence(
    'clean-bin',
    'check-fmt',
    'compile',
    'copy-binary',
    'package-binary',
    'package-fonts',
    'dist',
    'clean-home',
    'build-sample',
    'test',
    callback);
});

gulp.task('compile', function(callback) {
  //exec('go build timeline.go draw.go calendar.go data.go theme.go', function(err, stdout, stderr) {
  exec('go build', function(err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    callback(err);
  });
});

gulp.task('copy-binary', function() {
  return gulp.src(['./timeline', './timeline.exe'])
    .pipe(gulp.dest('../../../../bin'))
});

gulp.task('package-binary', function() {
  return gulp.src(['./timeline', './timeline.exe'], { base: '.' })
    .pipe(gulp.dest('package'))
});

gulp.task('package-fonts', function() {
  return gulp.src('./resource/font/*', { base: './resource/font' })
    .pipe(gulp.dest('package/resource/font'))
});

gulp.task('dist', function() {
  return gulp.src('./package/**/*', { base: './package' })
    .pipe(zip(pkg.name + '-' + pkg.version + '-' + platform + '.zip'))
    .pipe(gulp.dest('./dist'));
});

gulp.task('build-sample', function(callback) {
  exec('timeline data/sample.json data/sample_duration.json data/sample_unthemed.json data/sample_minimal.json', function(err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    callback(err);
  });
});

gulp.task('test', function(callback) {
  exec('go test', function(err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    callback(err);
  });
});

gulp.task('check-fmt', function(callback) {
  exec('gofmt -d timeline.go', function(err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    callback(err);
  });
});

gulp.task('clean-home', function() {
  return del.sync(['./timeline', './timeline.exe'], { force: true });
});

gulp.task('clean-bin', function() {
  return del.sync(['../../../../bin/timeline', './dist/' + pkg.name + '-' + pkg.version + '-' + platform + '.zip', './package/**/*'], { force: true });
});

gulp.task('watch', function() {
  gulp.watch(['./*.go', './data/*.json'], [
    'build'
  ]);
});