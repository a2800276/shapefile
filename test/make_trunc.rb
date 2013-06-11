
File.open("truncated.shp", "w") { |f|
  f.write(File.read("Geometrie_Wahlkreise_18DBT.shp")[0,20])
}
