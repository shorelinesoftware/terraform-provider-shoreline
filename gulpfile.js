const gulp = require('gulp');
const { dest, series, src } = require('gulp');
const replace = require('gulp-replace');
const rename = require('gulp-rename');
const shell = require('gulp-shell');
const data = require('./content/terms.json');

/**
 * Builds templates for the documentation.
 *
 * @returns {void}
 */
function buildTemplates() {
  return src('content/**/*.md')
    .pipe(
      replace(
        /\[(.+?)\]\(([^h].+?)\)/g,
        function handleReplace(match, p1, p2) {
          return `[${p1}](https://docs.shoreline.io${getTermUrl(p2)})`;
        }
      )
    )
    .pipe(
      rename(function (path) {
        // Updates the object in-place
        path.extname = '.md.tmpl';
      })
    )
    .pipe(dest('./'));
}

/**
 * Transforms a Docs term path to a URL.
 *
 * @param {string} value Term path to parse and transform
 * @returns Transformed URL, if applicable
 */
function getTermUrl(value) {
  let { terms } = data;

  // Ensure pattern isn't empty
  terms = terms.filter((term) =>
    term.patterns.filter((pattern) => pattern !== '')
  );
  // Process and filter term matches
  for (const term of terms) {
    for (const pattern of term.patterns) {
      const regex = new RegExp(`(^\/?t\\/${pattern}[s]?$)`);
      if (regex.exec(value)) {
        return `/${term.path}`;
      }
    }
  }

  const wildcardTerms = terms.filter((element) => element.wildcard);
  // Process remaining partial matches/wildcards
  for (const term of wildcardTerms) {
    for (const pattern of term.patterns) {
      const matches = value.match(
        new RegExp(`^(\/?t\/${pattern}[s]?)?(\.+?)?$`)
      );

      if (
        matches &&
        matches[1] &&
        matches[0].split('/').length - 1 === pattern.split('/').length + 1
      ) {
        if (matches[2] !== undefined) {
          // append extra wildcard path string
          return `/${term.path}${matches[2]}`;
        }
        return `/${term.path}`;
      }
    }
  }

  return value;
}

// Auto-generated documentation template files
gulp.task('generateDocs', shell.task('go generate'));

exports.buildTemplates = series(buildTemplates);
exports.build = series(buildTemplates, 'generateDocs');
exports.default = series(buildTemplates, 'generateDocs');
