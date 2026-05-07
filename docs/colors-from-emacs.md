# Colors from Emacs

`ntagcolor` supports CSS/web color tags by copying color names and hex values
from Emacs into the Go declarative color table.

The source expression is:

```elisp
(mapcar
 (lambda (c)
   (cons (substring-no-properties c)
         (get-text-property 0 'hex c)))
 (counsel-colors--web-alist))
```

`counsel-colors--web-alist` returns propertized strings. The visible string is
the lowercase color name, and the `hex` text property contains the normalized
`#rrggbb` value. Use `get-text-property` rather than `color-values`, because
`counsel-colors--web-alist` also carries compatibility fixes such as
`rebeccapurple -> #663399`.

To refresh the table:

```sh
emc-eval '(mapcar (lambda (c) (cons (substring-no-properties c) (get-text-property 0 (quote hex) c))) (counsel-colors--web-alist))'
```

Copy the resulting values into `tagStyles` as `Hex("#rrggbb")` entries. The Go
table also supports direct RGB declarations with `RGB(r, g, b)`.

Foreground handling:

- If `FG` is omitted, `ntagcolor` picks black or white text from background
  luminance.
- If `FG` is set with `Explicit(...)`, that foreground is used directly.

Precedence:

- Existing historical `ntagcolor` colors stay active when their names conflict
  with CSS colors.
- CSS alternatives for those conflicting names are kept as comments near the
  legacy entries.
- New CSS color names are added as normal declarative entries.
