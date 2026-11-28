module rename

sig Entry, Title {}
sig List {
  title: Entry -> lone Title
}

// Renaming one entry preserves every other entry and its title.
pred rename[prior, result: List, selected: Entry, replacement: Title] {
  some selected.(prior.title)
  result.title = prior.title ++ (selected -> replacement)
}

assert OtherTitlesSurvive {
  all a, b: List, selected, other: Entry, replacement: Title |
    rename[a, b, selected, replacement] and other != selected implies
      other.(a.title) = other.(b.title)
}

check OtherTitlesSurvive for 4
