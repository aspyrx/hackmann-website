module.exports = function(grunt) {
    grunt.initConfig({
        watch: {
            less: {
                files: ['src/style/*.less'],
                tasks: ['less:src']
            },
            go: {
                files: ['*.go'],
                tasks: ['go:build:hackmann']
            }
        },
        
        copy: {
            dist: {
                expand: true,
                cwd: 'src/',
                src: ['*.html', 'img/*'],
                dest: 'dist/'
            }
        },

        less: {
            src: {
                options: {
                    paths: 'src/'
                },
                files: {
                    'src/style/style.css': 'src/style/importer.less'
                }
            },
            dist: {
                options: {
                    paths: 'src/'
                },
                files: {
                    'src/style/style.css': 'src/style/importer.less'
                }
            }
        },
        
        cssmin: {
            dist: {
                files: [{
                    expand: true,
                    cwd: 'src/',
                    src: ['style/*.css'],
                    dest: 'dist/',
                    ext: '.css'
                }]
            }
        },
        
        uglify: {
            dist: {
                files: {
                    'dist/script/script.js': ['src/script/script.js']
                }
            }
        },
        
        go: {
            hackmann: {
                output: 'bin/hackmann-website.bin'
            }
        },
    });

    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-less');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-go');
    grunt.registerTask('build', ['copy:dist', 'less:dist', 'cssmin:dist', 'uglify:dist'])
    grunt.registerTask('default', ['watch']);
};
