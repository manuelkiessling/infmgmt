/*

- hier ist der befehlszeilen-zusammenbauer implementiert, der vom usecases layer benutzt wird um zB eine VM zu installieren
  oder puppet agent auszuführen usw.
- nutzt executor im infrastructure layer, um befehle tatsächlich auszuführen und deren ergebnis zu bekommen
- bekommt von aufrufer die infos als nackte daten, zB name der vm, größe arbeitsspeicher usw.
- webservice kann entscheiden einen json endpunkt anzubieten der alle machine infos auf einmal zurückgibt

*/
