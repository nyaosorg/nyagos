(defun md-to-html (c) (string-append "../docs/" (basename c) ".html"))

(defun make-title (md)
  (string-append
    "NYAOS.ORG/NYAGOS "
    (cond
      ((or (equal md "index.md") (equal md "index_ja.md"))
        "")
      ((or (match "_ja\.md$" md) (match "_en\.md$" md))
       (subseq md 0 (- (length md) 6)))
      (t
        (basename md)))))

(case $1
  (("clean")
   (dolist (html (wildcard "../docs/*.html"))
     (rm html)))
  (t
    (dolist (md (wildcard "*.md"))
      (if (not (match "^_" md))
        (let ((html (md-to-html md))
              (sidebar (if (match "_ja.md$" md) "_Sidebar_ja.md" "_Sidebar_en.md")))
          (if (updatep html md sidebar "_Header.md")
            (sh (format
                  nil
                  "minipage -sidebar \"~A\" -title \"~A\" _Header.md \"~A\" > \"~A\""
                  sidebar (make-title md) md html))
                  )))
      ) ; dolist
    ) ; t
  ) ; case
