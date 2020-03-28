'use strict'

const gulp = require('gulp')
const concatCss = require('gulp-concat-css')
const concat = require('gulp-concat')

const minify = require('gulp-minify')
const cleanCSS = require('gulp-clean-css')

const scripts = [
    './node_modules/jquery/dist/jquery.js',
    './node_modules/moment/moment.js',
    './js/instantclick.min.js',
    './js/jquery.poshytip.js',
    './js/jquery-editable-poshytip.js',
    './js/app.js']

const cssFiles = [
    './css/tachyons.min.css',
    './css/jquery-editable.css',
    './css/index.css',
]

gulp.task('scripts', function () {
    return gulp.src(scripts)
        .pipe(concat('bundle.js'))
        .pipe(gulp.dest('../static/'))
})

gulp.task('compress', function () {
    gulp.src('../static/bundle.js')
        .pipe(minify({
            ext: {
                src: '-debug.js',
                min: '.js'
            }
        }))
        .pipe(gulp.dest('../static'))
});

gulp.task('css', function () {
    return gulp.src(cssFiles)
        .pipe(concatCss('index.css'))
        .pipe(gulp.dest('../static/'))
})

gulp.task('minify-css', () => {
    return gulp.src('../static/index.css')
        .pipe(cleanCSS({compatibility: 'ie8'}))
        .pipe(gulp.dest('../static'));
});
