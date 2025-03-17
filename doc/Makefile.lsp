(defun md-to-html (c) (string-append "../docs/" (basename c) ".html"))

(case $1
  (("clean")
   (dolist (html (wildcard "../docs/*.html"))
     (rm html)))
  (t
    (dolist (md (wildcard "*.md"))
      (if (not (match "^_" md))
        (let ((html (md-to-html md))
              (sidebar (if (match "_ja.md$" md) "_Sidebar_ja.md" "_Sidebar_en.md")))
          (if (updatep html md sidebar "_Header.md" "wifky.css")
            (sh (string-append "minipage -header _Header.md -sidebar " sidebar " " md " > " html )))))
      ) ; dolist
    ) ; t
  ) ; case
