class OutlyerCli < Formula
  desc "Outlyer CLI allows to easily manage your Outlyer account via command line."
  homepage "https://www.outlyer.com/"
  url "https://github.com/outlyerapp/outlyer-cli/releases/download/1.0.0/outlyer-cli_1.0.0_Darwin_x86_64.tar.gz"
  version "1.0.0"
  sha256 "b9435f2909cadb6412faef9488cc26a6618cd7517862a28a1dce62a1cf34f102"

  def install
    bin.install "outlyer"
  end
end
